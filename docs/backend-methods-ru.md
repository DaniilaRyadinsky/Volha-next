# Методы backend Volha

Документ описывает сетевые методы backend и их связи с бизнес-логикой и хранилищами. Основа — текущая регистрация маршрутов в gateway, контракты gRPC из proto и реализации сервисов/репозиториев в каталоге backend.

## Как проходит запрос

1. Клиент вызывает HTTP-метод gateway (Echo).
2. Обработчик из internal/gateway/net/handlers проверяет HTTP-ввод и при необходимости cookie, параметры страниц или файлы.
3. Gateway вызывает одноимённый RPC-клиент Products или AuthService.
4. gRPC API преобразует protobuf-модель в доменную и вызывает продуктовый сервис либо auth-сервис.
5. Сервис проверяет бизнес-правила и обращается к репозиторию. Репозитории читают и изменяют PostgreSQL. Redis используется gateway только для кэша справочников, файлы хранятся локально в images/ и pdfs/.

Запросы проходят через middleware request-id и метрик/трейсов. Для обработчиков gateway срок ожидания gRPC-вызова обычно ограничен 3 секундами. Ошибки сервисов преобразуются в HTTP-статусы; основные соответствия: некорректный ввод — 400, конфликт уникальности — 409, отсутствие записи — 404, недоступный сервис — 503, тайм-аут — 504, внутренняя ошибка — 502.

### Доступ и общие правила

- Каталог, товары, статьи, справочники, слайды и просмотр фотографий доступны без входа. Middleware CheckId проверяет формат XID только если query-параметр id передан; конкретный обработчик требует его там, где он нужен.
- Изменяющие методы каталога, товаров, статей, слайдов, файлов и фотографий находятся под /api и включены только при режиме local/development/DEV. Нужен заголовок Authorization: Bearer <DEV_ADMIN_TOKEN>. В остальных режимах admin-маршруты возвращают 403.
- /api/auth/me, смена данных аккаунта и все методы /api/cart требуют access-cookie, которую gateway проверяет вызовом AuthService.Authenticate. Сами access и refresh JWT HttpOnly; endpoints регистрации/входа/обновления устанавливают их в cookie, но не возвращают в JSON.
- Идентификаторы доменных сущностей — XID. При создании товаров, статей и справочников ID создаёт gateway.
- Пагинация коллекций: page начинается с 1, по умолчанию 1; limit по умолчанию 20, максимум 100. В legacy-вызовах списка принимается offset вместо page, если page не задан.

## 1. Аккаунт и аутентификация

### HTTP API gateway

| HTTP-метод и путь | Обработчик | Что делает и с чем связан |
|---|---|---|
| POST /api/auth/register | Register | Принимает email, login, password; вызывает AuthService.Register; возвращает пользователя (201) и ставит access/refresh cookies. Сервис нормализует email/login, проверяет пароль, хеширует его bcrypt и создаёт пользователя вместе с первой сессией. |
| POST /api/auth/login | Login | Принимает identifier (логин или email) и пароль; вызывает AuthService.Login; при успехе выдаёт пару cookies и профиль пользователя. |
| POST /api/auth/refresh | Refresh | Берёт refresh-cookie, вызывает AuthService.Refresh, ротирует refresh-токен и устанавливает новую пару cookies. При отсутствующем или недействительном токене очищает cookies и отвечает 401. |
| POST /api/auth/logout | Logout | Если refresh-cookie присутствует, отзывает соответствующую сессию через AuthService.Logout; затем очищает обе cookies. Повторный вызов без cookie допустим. |
| GET /api/auth/me | Me | Требует авторизации middleware. Возвращает пользователя из principal, уже проверенного через AuthService.Authenticate. |
| PUT /api/auth/change-password | ChangePassword | Требует access-cookie; принимает old_password и new_password. Передаёт user_id и session_id из проверенного principal. После смены пароля текущая сессия сохраняется, остальные сессии пользователя отзываются. |
| PUT /api/auth/change-login | ChangeLogin | Требует access-cookie; принимает пароль и new_login. Проверяет пароль, обновляет логин в активной сессии; уникальность логина обеспечивает PostgreSQL. |

### Корзина

| HTTP-метод и путь | Обработчик | Функция и связи |
|---|---|---|
| GET /api/cart | GetCart | Возвращает корзину текущего пользователя. В каждой строке присоединяет продукт: slug, title, обычную/акционную цену и первую фотографию. |
| POST /api/cart/items | AddCartItem | Добавляет product_id и quantity. Если товар уже есть, количество увеличивается; допустимый итог — не больше 99. Возвращает обновлённую корзину. |
| PUT /api/cart/items/{product_id} | UpdateCartItem | Устанавливает новое количество существующей позиции (1–99), затем возвращает всю корзину. product_id берётся из пути. |
| DELETE /api/cart/items/{product_id} | RemoveCartItem | Удаляет позицию по пользователю и product_id; возвращает новую корзину. |
| DELETE /api/cart | ClearCart | Удаляет все позиции пользователя и возвращает пустую корзину. |

Все HTTP-методы корзины защищены RequireAuth. user_id нельзя передать через тело запроса: gateway берёт его из аутентифицированной cookie-сессии. Auth-сервис хранит строки в cart_items, а GetCart считает line_total по акционной цене, если она ненулевая, иначе по обычной; итоговая сумма и количество суммируются по позициям. Количество добавления/обновления проверяется сервисом, а добавление дополнительно проверяет существование товара.

### Внутренние RPC AuthService и цепочка хранения

| RPC / метод Service | Назначение | Слой данных |
|---|---|---|
| Register | Проверить и нормализовать учётные данные, создать пользователя и сессию, выдать токены. Email: не более 254 символов и корректный адрес; login: 3–32 символа, латиница/цифры/точка/дефис/подчёркивание; пароль: 10–72 байта, буква и цифра. | Store.CreateUserAndSession — транзакционная запись users + auth_sessions. |
| Login | Найти пользователя по email или login без учёта регистра, сравнить bcrypt-хеш, создать новую сессию и выдать токены. | FindUserForLogin → CreateSession. |
| Refresh | Проверить подпись, срок и тип refresh JWT; заменить сохранённый хеш refresh-токена для активной сессии. Повторное использование старого токена не проходит. | RotateSession блокирует запись сессии и атомарно меняет хеш/срок. |
| Logout | Отозвать refresh-сессию по user/session и хешу токена. Невалидный токен не вызывает ошибку, что делает операцию идемпотентной. | RevokeSession удаляет auth_sessions. |
| Authenticate | Проверить access JWT и наличие незавершённой сессии; вернуть пользователя и session_id для gateway middleware. | AuthenticateSession читает users через активную auth_sessions. |
| ChangePassword | Проверить новый пароль, старый пароль и запрет повторного текущего пароля; сохранить новый bcrypt-хеш. | PasswordHash → ChangePassword. Репозиторий блокирует текущую сессию, меняет users.password_hash и отзывает все другие сессии. |
| ChangeLogin | Нормализовать login, проверить пароль и обновить логин при активной сессии. | PasswordHash → ChangeLogin; активная сессия блокируется транзакционно. Конфликт уникальности отображается как AlreadyExists. |
| GetCart | Прочитать корзину с данными продуктов и агрегатами. | Store.GetCart, таблицы cart_items + products. |
| AddCartItem | Проверить XID продукта и quantity 1–99; добавить строку или увеличить имеющуюся, если сумма не превышает 99; вернуть корзину. | Store.AddCartItem → GetCart. Несуществующий продукт — NotFound, превышение — ResourceExhausted. |
| UpdateCartItem | Проверить продукт и количество 1–99; заменить количество существующей строки; вернуть корзину. | Store.UpdateCartItem → GetCart. Нет строки — NotFound. |
| RemoveCartItem | Проверить XID, удалить строку; вернуть корзину. | Store.RemoveCartItem → GetCart. Нет строки — NotFound. |
| ClearCart | Очистить все строки пользователя; вернуть корзину. | Store.ClearCart → GetCart. |

Дополнительные методы auth-хранилища: New создаёт Store; lockActiveSession блокирует активную сессию перед изменением аккаунта; mapPostgresError переводит SQL-ограничения в доменные ошибки. JWT выпускаются методами issuePair/signToken и проверяются verifyToken; в БД сохраняется SHA-256 хеш refresh-токена, а не сам токен.

## 2. Каталог товаров и поиск

| HTTP-метод и путь | Обработчик | RPC → сервис → репозиторий | Функция |
|---|---|---|---|
| POST /api/product | CreateProduct | Products.CreateProduct → ProductsService.CreateProduct → ProductRepo.CreateProduct | Создаёт товар. Gateway генерирует XID, проверяет, что все имена фотографий есть на диске. |
| PUT /api/product | UpdateProduct | Products.UpdateProduct → UpdateProduct → ProductRepo.UpdateProduct | Заменяет переданные атрибуты товара. Перед записью читает старые фотографии/фото по цветам; после обновления инвалидирует справочники и удаляет только больше нигде не используемые файлы. |
| DELETE /api/product?id=… | DeleteProduct | Products.DeleteProduct → DeleteProduct → ProductRepo.DeleteProduct | Сначала загружает товар и его фотографии по цветам, затем удаляет товар. После успеха сбрасывает кэш справочников и очищает неиспользуемые изображения. |
| GET /api/product?id=… | GetProduct | Products.GetProduct → GetProduct → ProductRepo.GetProduct | Возвращает подробную карточку товара с брендом, категорией, страной, материалами, цветами, фотографиями и похожими товарами. |
| GET /api/product/all?page=…&limit=… | GetAllProducts | Products.GetAllProducts → GetAllProducts → ProductRepo.GetAllProducts | Возвращает постраничный список карточек каталога, total, page и limit. Список сортируется по названию и ID. |
| POST /api/product/filter | FilterProducts | Products.FilterProducts → FilterProducts → ProductRepo.FilterProducts | Фильтрует и сортирует каталог, возвращает компактные ProductCatalogItem и общее число совпадений. |

Создание и изменение принимают ProductId: id, title, article, brand, category, country, width/height/depth, materials, colors, photos, seems, price, description, isFeatured, pricesale и slug. Проверки сервиса требуют XID для сущности и ссылок, непустые title/slug/article, артикул ровно из 8 цифр, цену > 0, скидочную цену ≥ 0, размеры > 0, хотя бы одно фото; списки материалов/цветов/похожих ID не должны содержать повторы, товар не может ссылаться на себя как похожий.

Фильтр принимает brand/category/country, материалы/цвета, диапазоны ширины/высоты/глубины/цены, sort_by, sort_order, page/limit и текст поиска. Разрешённая сортировка: price, width, height, depth, title; направление asc/desc. Текст search gateway переводит в поле title; SQL использует подстрочный поиск ILIKE. При фильтрации по материалам/цветам товар подходит, если связь найдена хотя бы с одним из переданных ID. limit в сервисе не больше 100.

### Связи товарного каталога

- products.brand_id, category_id и country_id ссылаются на справочные таблицы brands, categories и countries.
- Материалы и цвета связаны с товарами через product_materials и product_colors; эти же данные используются в карточке и фильтрах.
- Поле seems хранит связи похожих товаров. ProductRepo собирает их в карточке вместе с краткими данными товара.
- product_color_photos связывает пару product_id + color_id с набором изображений. Эти изображения включаются в проверки ссылок на файлы и очищаются вместе с заменой/удалением товара.
- cart_items.product_id связан с products.id; карточка корзины вычисляет цену и изображение из товара.

## 3. Справочники и контент каталога

У каждого ресурса ниже есть пять отдельных HTTP/gRPC методов. HTTP-обработчики вызывают одноимённый RPC Products, тот — одноимённый метод ProductsService, далее соответствующий CatalogRepo-метод в PostgreSQL. Create генерирует ID XID в gateway; чтение одной записи и удаление используют query id; update получает объект целиком с id в JSON. Сервис проверяет XID, обязательные поля и выполняет create/update/delete внутри транзакции. GetAll сортирует записи по названию и ID, кроме слайдов, которые сортируются по ID. Ошибки внешних ключей не позволяют удалить запись, используемую другими таблицами.

| Ресурс / HTTP путь | Метод: HTTP → handler → RPC/Service → Repo | Тело / результат и назначение |
|---|---|---|
| Бренды /api/brand | POST CreateBrand → CreateBrand → CreateBrand → CreateBrand | Тело title; создаёт бренд и сбрасывает кэш справочников; ответ содержит новый id. |
| Бренды /api/brand | PUT UpdateBrand → UpdateBrand → UpdateBrand → UpdateBrand | Тело id,title; переименовывает бренд и сбрасывает кэш. |
| Бренды /api/brand?id=… | DELETE DeleteBrand → DeleteBrand → DeleteBrand → DeleteBrand | Удаляет бренд; внешний ключ запрещает удаление, если он привязан к товару. Сбрасывает кэш. |
| Бренды /api/brand/all | GET GetAllBrands → GetAllBrands → GetAllBrands → GetAllBrands | Возвращает весь список брендов. |
| Бренды /api/brand?id=… | GET GetBrand → GetBrand → GetBrand → GetBrand | Возвращает один бренд по XID. |
| Категории /api/category | POST CreateCategory → CreateCategory → CreateCategory → CreateCategory | Тело slug,title,description,img; ID генерируется; проверяет наличие файла img; сбрасывает кэш. |
| Категории /api/category | PUT UpdateCategory → UpdateCategory → UpdateCategory → UpdateCategory | Тело id,slug,title,description,img; проверяет новый img, сбрасывает кэш и при необходимости удаляет старую картинку, если на неё больше нет ссылок. |
| Категории /api/category?id=… | DELETE DeleteCategory → DeleteCategory → DeleteCategory → DeleteCategory | Удаляет категорию; вернёт ошибку внешнего ключа, если к ней привязаны товары; инвалидирует кэш и очищает больше не используемую картинку. |
| Категории /api/category/all | GET GetAllCategories → GetAllCategories → GetAllCategories → GetAllCategories | Возвращает категории вместе с числом товаров в каждой. |
| Категории /api/category?id=… | GET GetCategory → GetCategory → GetCategory → GetCategory | Возвращает категорию и вычисляемый count товаров. |
| Страны /api/country | POST CreateCountry → CreateCountry → CreateCountry → CreateCountry | Тело title,friendly; генерирует ID и сбрасывает кэш. |
| Страны /api/country | PUT UpdateCountry → UpdateCountry → UpdateCountry → UpdateCountry | Тело id,title,friendly; изменяет страну и сбрасывает кэш. |
| Страны /api/country?id=… | DELETE DeleteCountry → DeleteCountry → DeleteCountry → DeleteCountry | Удаляет страну, если её не блокируют внешние ключи; инвалидирует кэш. |
| Страны /api/country/all | GET GetAllCountries → GetAllCountries → GetAllCountries → GetAllCountries | Возвращает список стран. |
| Страны /api/country?id=… | GET GetCountry → GetCountry → GetCountry → GetCountry | Возвращает одну страну по ID. |
| Материалы /api/material | POST CreateMaterial → CreateMaterial → CreateMaterial → CreateMaterial | Тело title; создаёт материал и инвалидирует справочники. |
| Материалы /api/material | PUT UpdateMaterial → UpdateMaterial → UpdateMaterial → UpdateMaterial | Тело id,title; изменяет материал и инвалидирует справочники. |
| Материалы /api/material?id=… | DELETE DeleteMaterial → DeleteMaterial → DeleteMaterial → DeleteMaterial | Удаляет материал, если его не удерживает внешний ключ; инвалидирует кэш. |
| Материалы /api/material/all | GET GetAllMaterials → GetAllMaterials → GetAllMaterials → GetAllMaterials | Возвращает список материалов. |
| Материалы /api/material?id=… | GET GetMaterial → GetMaterial → GetMaterial → GetMaterial | Возвращает один материал по ID. |
| Цвета /api/color | POST CreateColor → CreateColor → CreateColor → CreateColor | Тело title,hex; генерирует ID, сбрасывает кэш. |
| Цвета /api/color | PUT UpdateColor → UpdateColor → UpdateColor → UpdateColor | Тело id,title,hex; обновляет цвет и сбрасывает кэш. |
| Цвета /api/color?id=… | DELETE DeleteColor → DeleteColor → DeleteColor → DeleteColor | Удаляет цвет, если нет зависимостей по внешнему ключу; инвалидирует справочники. |
| Цвета /api/color/all | GET GetAllColors → GetAllColors → GetAllColors → GetAllColors | Возвращает список цветов. |
| Цвета /api/color?id=… | GET GetColor → GetColor → GetColor → GetColor | Возвращает один цвет по ID. |

## 4. Слайды главной страницы

Slides — отдельный контентный каталог, не часть фильтруемых справочников. Структура объекта: id, link, img, img762, title, description, badge. CRUD-методы следуют общей цепочке Products RPC → ProductsService → CatalogRepo.

| HTTP-метод и путь | Обработчик / RPC / сервис | Функция и связи |
|---|---|---|
| POST /api/slide | CreateSlide / CreateSlide / CreateSlide | Создаёт слайд, ID генерируется; проверяет наличие файлов img и img762; возвращает ID. |
| PUT /api/slide | UpdateSlide / UpdateSlide / UpdateSlide | Обновляет поля слайда, проверяет новые изображения; старые файлы удаляет только если на них больше нет ссылок. |
| DELETE /api/slide?id=… | DeleteSlide / DeleteSlide / DeleteSlide | Читает слайд, удаляет запись, затем очищает неиспользуемые img и img762. |
| GET /api/slide/all | GetAllSlides / GetAllSlides / GetAllSlides | Возвращает все слайды. |
| GET /api/slide?id=… | GetSlide / GetSlide / GetSlide | Возвращает один слайд по ID. |

Репозиторий предоставляет методы GetAllSlides, GetSlide, CreateSlide, UpdateSlide, DeleteSlide. Сервис валидирует XID, link, img и img762; gateway дополнительно удостоверяется, что оба изображения существуют.

## 5. Статьи и PDF

Статья имеет id, slug, title, img[] и text (ссылка на PDF); creation_time возвращается в миллисекундах. Фактический PDF хранится в pdfs/, картинки — в images/.

| HTTP-метод и путь | Handler → RPC → сервис → репозиторий | Назначение |
|---|---|---|
| POST /api/article | CreateArticle → CreateArticle → CreateArticle → CreateArticle | Создаёт статью; gateway генерирует ID, нормализует text до /pdfs/<имя>.pdf и проверяет наличие PDF и всех img. Время создания ставит БД. |
| PUT /api/article | UpdateArticle → UpdateArticle → UpdateArticle → UpdateArticle | Изменяет статью, оставляет прежнее creation_time; после обновления очищает снятые изображения и PDF, если они нигде больше не используются. |
| DELETE /api/article?id=… | DeleteArticle → DeleteArticle → DeleteArticle → DeleteArticle | Удаляет статью, затем проверяет ссылки и очищает её изображения и PDF. |
| GET /api/article?id=… | GetArticle → GetArticle → GetArticle → GetArticle | Возвращает статью по XID. |
| GET /api/article/all?page=…&limit=… | GetAllArticles → GetAllArticles → GetAllArticles → GetAllArticles | Возвращает страницу статей, total, page, limit; порядок — новые сначала. |

Service.validateArticle требует XID, title и slug, путь /pdfs/ и простое имя PDF с расширением .pdf; для каждой картинки проверяет, что задано только имя файла. Gateway bindArticle также проверяет существование PDF и изображений в файловой системе. Репозиторий: GetAllArticles, GetArticle, CreateArticle, UpdateArticle, DeleteArticle; создание записывает время, обновление не меняет created_at.

## 6. Фотографии товара в разрезе цвета

Связь хранит product_id, color_id и массив photos. В БД для пары товара и цвета предусмотрена отдельная запись; product/color должны существовать, повторное создание такой пары конфликтует.

| HTTP-метод и путь | Handler → RPC → Service → Repo | Назначение |
|---|---|---|
| POST /api/colorphotos | CreateProductColorPhotos → CreateProductColorPhotos → CreateProductColorPhotos → CreateProductColorPhotos | Создаёт набор фотографий для пары product/color; gateway проверяет файлы, сервис требует два XID и непустой список фотографий. |
| PUT /api/colorphotos | UpdateProductColorPhotos → UpdateProductColorPhotos → UpdateProductColorPhotos → UpdateProductColorPhotos | Заменяет список фотографий; после записи очищает только снятые и нигде больше не используемые файлы. |
| DELETE /api/colorphotos | DeleteProductColorPhotos → DeleteProductColorPhotos → DeleteProductColorPhotos → DeleteProductColorPhotos | Тело ProductColorPhotosId {product_id,color_id}; сначала считывает старые имена, удаляет связь, очищает неиспользуемые файлы. |
| GET /api/colorphotos/all | GetAllProductColorPhotos → GetAllProductColorPhotos → GetAllProductColorPhotos → GetAllProductColorPhotos | Возвращает все записи связи. |
| POST /api/colorphotos/photos/get | GetPhotosByProductAndColor → GetPhotosByProductAndColor → GetPhotosByProductColor → GetPhotosByProductColor | Принимает пару ID в JSON и возвращает массив фото именно для неё. |

Дополнительные внутренние RPC Products.GetProductColorPhotosByProduct (id товара) возвращает все фото по цветам этого товара. Отдельного HTTP-маршрута для этого RPC нет; его использует обработчик обновления/удаления товара для сбора файлов на очистку.

## 7. Справочники для фильтров и Redis-кэш

| HTTP-метод и путь | Обработчик → RPC → сервис → репозиторий | Назначение |
|---|---|---|
| GET /api/dictionaries/all | GetAllDictionaries → GetDictionaries → GetDictionaries → DictRepo.GetDictionaries | Возвращает бренды, категории, страны, материалы, цвета и min/max цены и размеров для общего каталога. При отсутствии Redis-кэша запрашивает product-service и сохраняет результат. |
| GET /api/dictionaries/all/by-category?id=… | GetAllDictionariesByCategory → GetDictionaries → GetDictionaries → DictRepo.GetDictionariesByCategoryID | Возвращает те же справочные данные и агрегаты для категории. Сервис сначала проверяет, что категория существует. Результат кэшируется отдельным ключом категории. |

Специальный ID CLR обозначает общий набор справочников без фильтра по категории. Redis использует поколение кэша в ключах; CleanDictionaries увеличивает поколение вместо удаления каждого ключа. Gateway инвалидирует поколение при успешном изменении товара или справочников бренда/категории/страны/материала/цвета. Изменения статей, слайдов и фото по цветам кэш справочников не сбрасывают. При ошибке Redis gateway пишет её в лог и пытается получить данные из product-service.

## 8. Загрузка, раздача и удаление файлов

| HTTP-метод и путь | Метод | Назначение и связи |
|---|---|---|
| POST /api/files | UploadFile | Принимает multipart-поле file или img. По содержимому разрешает JPEG, PNG, WEBP, GIF и PDF; имя формируется UUID, предел — 10 МБ. Файлы попадают в images/ либо pdfs/ по расширению; возвращается name. |
| DELETE /api/files?title=… | DeleteFile | Получает только имя файла. Через Products.IsFileReferenced проверяет ссылки в категориях, товарах, фото товара/цвета, слайдах и статьях. При наличии ссылки отвечает 409; иначе удаляет из локального хранилища. |
| GET /images/* | Static middleware | Раздаёт загруженные изображения из images/. |
| GET /pdfs/* | Static middleware | Раздаёт загруженные PDF из pdfs/. |

Удаление и автоочистка используют imageFilePath/pdfFilePath и validateImageFilename для запрета путей с каталогами/разделителями. ensureImageReferences/ensurePDFReference не дают сохранить запись со ссылкой на отсутствующий файл. cleanupFileIfUnreferenced повторно сверяется с БД перед удалением, поэтому общее изображение не удаляется, пока на него ссылается другая запись. RPC IsFileReferenced валидирует простое имя и передаёт запрос FileReferenceRepo.IsFileReferenced.

## 9. Полный внутренний контракт Products

Все методы этого контракта определены в api/proto/filez/products.proto и реализованы в internal/product-service/api. Для обычных CRUD выше уже приведены HTTP-маршрут, функция и цепочка. Здесь перечислены RPC-имена целиком, включая методы без HTTP-фасада.

| RPC Products | Для чего служит | Сервис и репозиторий |
|---|---|---|
| CreateArticle, UpdateArticle, DeleteArticle, GetArticle, GetAllArticles | CRUD статей и постраничный список. | ProductsService CRUD статей → ArticleRepo методы с теми же именами. |
| GetDictionaries | Получить полный или категорийный набор фильтров; CLR/пустой ID — полный набор. | ProductsService.GetDictionaries → DictRepo.GetDictionaries или GetDictionariesByCategoryID; для category ID также CatalogRepo.GetCategory. |
| CreateProduct, UpdateProduct, DeleteProduct, GetAllProducts, GetProduct, FilterProducts | Создать/изменить/удалить товар, получить карточку/страницу или каталог по фильтру. | ProductsService методы с теми же именами → ProductRepo методы с теми же именами. |
| CreateBrand, UpdateBrand, DeleteBrand, GetAllBrands, GetBrand | CRUD брендов. | ProductsService → CatalogRepo методы бренда. |
| CreateCategory, UpdateCategory, DeleteCategory, GetAllCategories, GetCategory | CRUD категорий с вычисляемым количеством товаров. | ProductsService → CatalogRepo методы категории. |
| CreateCountry, UpdateCountry, DeleteCountry, GetAllCountries, GetCountry | CRUD стран происхождения товара. | ProductsService → CatalogRepo методы страны. |
| CreateMaterial, UpdateMaterial, DeleteMaterial, GetAllMaterials, GetMaterial | CRUD материалов и справочных связей продуктов. | ProductsService → CatalogRepo методы материала. |
| CreateColor, UpdateColor, DeleteColor, GetAllColors, GetColor | CRUD цветов и связей product_colors/product_color_photos. | ProductsService → CatalogRepo методы цвета. |
| CreateSlide, UpdateSlide, DeleteSlide, GetAllSlides, GetSlide | CRUD слайдов. | ProductsService → CatalogRepo методы слайдов. |
| CreateProductColorPhotos, UpdateProductColorPhotos, DeleteProductColorPhotos | Управление набором фото для пары product/color. | ProductsService → ProductColorPhotosRepo методы с теми же именами. |
| GetAllProductColorPhotos, GetProductColorPhotosByProduct, GetPhotosByProductAndColor | Получение всех связей, связей по товару и фото для одной пары товара/цвета. | ProductsService.GetAllProductColorPhotos, GetProductColorPhotosByProduct, GetPhotosByProductColor → соответствующие ProductColorPhotosRepo методы. |
| IsFileReferenced | Проверяет использование имени файла в таблицах контента. | ProductsService.IsFileReferenced → FileReferenceRepo.IsFileReferenced. |

В gRPC CRUD-обёртках handleCRUDResponse преобразует доменные ошибки в gRPC status: InvalidArgument, NotFound, AlreadyExists, FailedPrecondition, Unimplemented или Internal. Преобразования protobuf ↔ domain выполняются функциями *FromRpc/*ToRpc в пакетах product-service/domain.

## 10. Поддерживающие методы сервера

- GET /health/live → health.Live возвращает состояние процесса без проверки зависимостей.
- GET /health/ready проверяет Ready обоих gRPC-клиентов gateway; каждый Ready проверяет состояние соединения и вызывает gRPC health.Check. Если один из сервисов не готов, endpoint возвращает ошибку готовности.
- gRPC health.Check зарегистрирован на product-service и auth-service. Product dependencyHealthServer.Check и auth healthServer.Check вызывают переданную проверку БД с тайм-аутом 2 секунды и отвечают SERVING/NOT_SERVING.
- GET /swagger/* обслуживается swagger UI из backend/docs и описывает HTTP-API.
- Gateway: echo.New собирает middleware и маршруты; Echo.Run запускает HTTP-сервер; Echo.Stop выполняет graceful shutdown; MustRun логирует ошибку запуска. skipObservabilityRoute исключает health/metrics/pprof из части telemetry.
- gRPC-клиенты gateway: products.New/auth.New устанавливают соединение и ждут Ready до 30 секунд; Ready дополнительно вызывает health.Check с timeout конфигурации; Close закрывает соединение. Auth-клиент также предоставляет ValidateAccessToken, который вызывает AuthService.Authenticate.
- Product transactions: ProductsService.runInTx запускает изменяющие операции через TxRunner.RunInTx; RepoProvider выбирает Product/Catalog/Dict/Pcp/FileRefs/Article репозиторий с учётом SQL-транзакции в context. TxFromContext достаёт транзакцию для SQL-операций.
- PostgreSQL подключения: product-service.NewConnect и auth-service.Connect открывают/проверяют DB; Ping используется проверками готовности; Disconnect/Close освобождают подключение.
- Команды обслуживания: cmd/product-service/migrator выполняет миграции up/down; cmd/dumper загружает настройки, делает backup через pg_dump и чистит старые копии. Это CLI-методы, а не HTTP/gRPC API.

## Основные файлы исходников

- Маршруты HTTP: backend/internal/gateway/net/echo/setup.go.
- HTTP обработчики: backend/internal/gateway/net/handlers/.
- Proto-контракты: backend/api/proto/filez/products.proto, auth.proto и domain.proto.
- gRPC обработчики: backend/internal/product-service/api/ и backend/internal/auth-service/api/.
- Бизнес-логика: backend/internal/product-service/service/ и backend/internal/auth-service/service/.
- PostgreSQL: backend/internal/product-service/repository/ и backend/internal/auth-service/repository/.
