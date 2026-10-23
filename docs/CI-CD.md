# CI/CD для Volha: текущий этап CI

Этот этап проверяет коммит без развёртывания приложения. Порядок Jenkins pipeline:
сборка Go и frontend → модульные Go-тесты → `go vet` и SonarQube → проверка
зависимостей → проверка безопасности исходников → сборка и проверка Docker-образов.
Каждая команда завершает stage с ошибкой при нарушении своего порога. Следующие
stage не запускаются; уже полученные отчёты остаются в артефактах сборки.

Интеграционные тесты с PostgreSQL/Redis, Testcontainers, E2E, нагрузочные тесты
и развёртывание **не входят** в этот pipeline. Их подключим на отдельном этапе
CDL. Исходники frontend не изменяются: здесь выполняются только установка
зависимостей по lockfile и сборка web/admin.

## Где лежат настройки

| Файл | Назначение |
| --- | --- |
| `.github/workflows/backend-ci.yml` | Быстрые Go-проверки при push в `main` и pull request; интеграционные тесты из него убраны. |
| `.github/workflows/frontend-ci.yml` | Сборка двух frontend-приложений при push в `main`. |
| `.github/workflows/security.yml` | Gitleaks, Semgrep CE и Trivy при push в `main`; отчёты в артефактах Actions. |
| `ci/jenkins/compose.yaml` | Локальные Jenkins, Docker-in-Docker, SonarQube, PostgreSQL для SonarQube и Nexus. |
| `ci/jenkins/Dockerfile` | Jenkins с Docker CLI и плагинами Pipeline, Git, JUnit, Credentials Binding. |
| `ci/jenkins/Jenkinsfile` | Порядок stage, учётные данные, публикация отчётов, quality gates. |
| `ci/jenkins/run.sh` | Точные команды сборок и сканеров; версии инструментов закреплены здесь. |
| `sonar-project.properties` | Область анализа Go-кода и импорт `ci-reports/coverage.out`. |
| `ci/jenkins/.env.example` | Образец локального пароля БД SonarQube. |

Генерируемые `ci-reports/` и `.ci-cache/` игнорируются Git. Jenkins хранит
состояние в Docker volumes, поэтому `docker compose down` его сохраняет.

## Как поднять локальный CI

Нужны Docker Desktop в режиме Linux containers, Docker Compose и достаточно
памяти для одновременно работающих Jenkins, SonarQube, PostgreSQL и Nexus
(практически стоит выделить Docker не менее 10–12 ГБ). Нужен исходящий доступ
к Go/npm registry, Semgrep Registry, базам OSV/Go/Trivy и Docker Hub.
Команды ниже выполняются
из корня репозитория в PowerShell:

```powershell
Copy-Item ci/jenkins/.env.example ci/jenkins/.env
# Откройте ci/jenkins/.env и замените SONAR_DB_PASSWORD на свой длинный пароль.
docker compose -f ci/jenkins/compose.yaml --env-file ci/jenkins/.env config --quiet
docker compose -f ci/jenkins/compose.yaml --env-file ci/jenkins/.env up -d --build
docker compose -f ci/jenkins/compose.yaml --env-file ci/jenkins/.env ps
```

Адреса интерфейсов на компьютере: Jenkins — `http://localhost:8085`, SonarQube —
`http://localhost:9000`, Nexus — `http://localhost:8086`. Порты привязаны только
к `127.0.0.1`. Nexus Docker registry через HTTP connector доступен из Jenkins
как `nexus:8082`, а с компьютера как `localhost:8087`. До настройки репозитория
в Nexus порт 8082 ничего не обслуживает.

Первоначальный пароль Jenkins покажет команда:

```powershell
docker compose -f ci/jenkins/compose.yaml --env-file ci/jenkins/.env exec jenkins cat /var/jenkins_home/secrets/initialAdminPassword
```

Завершите мастер настройки Jenkins. Плагины из Dockerfile уже установлены.
SonarQube при первой установке использует `admin` / `admin`; сразу измените
пароль. Начальный пароль Nexus находится в `/nexus-data/admin.password`:

```powershell
docker compose -f ci/jenkins/compose.yaml --env-file ci/jenkins/.env exec nexus cat /nexus-data/admin.password
```

Смените пароль Nexus при первом входе. Не сохраняйте эти пароли или `.env` в Git.
Остановка без удаления данных:

```powershell
docker compose -f ci/jenkins/compose.yaml --env-file ci/jenkins/.env down
```

### Настройка SonarQube

1. В SonarQube создайте локальный проект с ключом `volha-backend` и токен
   анализа. В Jenkins откройте **Manage Jenkins → Credentials → System →
   Global credentials → Add Credentials**, тип **Secret text**, ID `sonar-token`;
   вставьте токен. Pipeline передаёт его сканеру через переменную окружения.
2. У проекта выберите Quality Gate. Можно начать со встроенного **Sonar way**;
   для строгого режима создайте свой gate с требованием рейтингов **A** для
   безопасности и надёжности нового кода и покрытием нового кода не ниже 80%.
   Порог покрытия текущего кода задавайте после просмотра стартового отчёта.
   Точные условия gate управляются в интерфейсе SonarQube, не в `Jenkinsfile`.
3. В `sonar-project.properties` указаны анализируемые `backend/cmd` и
   `backend/internal`, исключены сгенерированные файлы и Go-тесты из исходников.
   Тесты учитываются отдельно, покрытие импортируется из Go `coverprofile`.
   `sonar.qualitygate.wait=true` заставляет сканер дождаться статуса gate
   (до 600 секунд) и завершиться с ошибкой, если он не пройден. Webhook не нужен.

### Настройка Jenkins job и запуск

Создайте **New Item → Pipeline**. В разделе **Pipeline** выберите **Pipeline
script from SCM**, SCM **Git**, укажите URL этого репозитория, ветку `*/main`,
Script Path `ci/jenkins/Jenkinsfile`. Для закрытого репозитория добавьте
учётные данные Git в Jenkins. Используйте полный clone: Gitleaks проверяет
историю коммитов, доступную в checkout.

Для ручного запуска нажмите **Build Now**. Чтобы Jenkins на локальном компьютере
замечал push в `main` без публичного webhook, включите **Poll SCM** с расписанием
`H/5 * * * *`: Jenkins проверит ветку примерно раз в пять минут и запустит job
при новом коммите. Локальный `git commit` без push Jenkins по этому сценарию
не запускает. Для проверки текущего коммита из интерфейса используйте
**Build with Parameters** и оставьте `PUBLISH_IMAGES=false`.

Важная деталь Docker-in-Docker: оба контейнера используют один том
`jenkins-data` по пути `/var/jenkins_home`. Поэтому scanner-контейнеры видят
workspace Jenkins через bind mount. Docker daemon защищён TLS и не опубликован
на хосте. Режим `privileged` нужен для локального dind. Для будущего удалённого
Jenkins потребуется отдельный защищённый агент и registry по HTTPS.

## Проверки, отчёты и quality gates

| Этап | Инструмент и конфигурация | Что считается ошибкой | Где смотреть результат |
| --- | --- | --- | --- |
| Сборки | Go 1.25.4: `go mod verify`, `go build ./...`; Node 24 + pnpm 11.3.0: `pnpm install --frozen-lockfile`, `pnpm build` | Ошибка зависимостей или сборки | Jenkins Stage View и Console Output |
| Модульные тесты | `gotestsum` 1.13.0 запускает `go test -race -count=1 -covermode=atomic` без тега `integration` | Падение теста или гонка | Jenkins **Test Result** из `ci-reports/unit.xml`; профиль `coverage.out` |
| Статический анализ | `go vet`; SonarScanner CLI 11 и `sonar-project.properties` | Ошибка `vet`, анализа или Quality Gate SonarQube | Jenkins Console Output; проект `volha-backend` в SonarQube |
| Go-зависимости | `govulncheck` 1.7.0 анализирует доступные уязвимые вызовы | Ненулевой код инструмента, в том числе найденная затрагивающая код уязвимость | `ci-reports/govulncheck.txt` |
| Open-source зависимости | OSV-Scanner 2.3.2 анализирует `backend/go.mod` | Любая известная уязвимость в Go-зависимостях; это намеренно строже порога CRITICAL | `ci-reports/osv.txt` |
| Безопасность Go-кода | gosec 2.29.0 по backend, исключены сгенерированные файлы | Находка с серьёзностью HIGH и уверенностью MEDIUM+ или ошибка анализа | `ci-reports/gosec.txt` |
| Правила SAST | Semgrep CE 1.182.0, набор `p/golang` | Находка уровня ERROR или ошибка сканера | `ci-reports/semgrep.sarif` |
| Секреты | Gitleaks 8.30.1 проверяет доступную историю Git, значение секрета маскируется | Любая находка | Jenkins Console Output; `ci-reports/gitleaks.json` |
| Образы и Dockerfile | Multi-stage `backend/Dockerfile` собирает 5 целей; Trivy 0.76.1 сканирует образы и конфигурацию backend | HIGH/CRITICAL уязвимость или ошибка конфигурации | `ci-reports/trivy-*.txt` |

Чтобы повторить отдельную проверку без Jenkins, используйте команды из
`ci/jenkins/run.sh`: каждая ветка `case` соответствует имени stage. Например,
из корня репозитория можно запустить `go -C backend test ./...`,
`go -C backend vet ./...` и проверку Go lockfile:

```powershell
docker run --rm --volume "${PWD}:/src" --workdir /src ghcr.io/google/osv-scanner:v2.3.2 scan source --lockfile=backend/go.mod --format=table
```

Другие scanner-контейнеры вызываются тем же способом, с версией и параметрами
из `run.sh`. Для SonarQube нужен работающий локальный сервер и токен, для Trivy
образов — предварительно собранные Docker-образы.

После каждого запуска откройте job → номер сборки → **Stage View**, **Test
Result**, **Console Output** и **Artifacts**. Все созданные `ci-reports/**`
архивируются даже при падении stage. Отчёты JSON/SARIF предназначены для
загрузки в просмотрщики; текстовые отчёты OSV, govulncheck и Trivy читаются
непосредственно. SonarQube показывает тренд по новому коду и детализацию
нарушений. Нулевой exit code сканера обязателен: ошибки загрузки правил,
базы уязвимостей или сети не превращаются в успешную проверку.

На 23 октября 2026 года первая проверка текущих зависимостей уже даёт красный
результат: OSV-Scanner нашёл 42 записи по трём Go-пакетам, включая две HIGH;
`govulncheck` обнаружил вызываемые уязвимые функции, а Trivy нашёл
HIGH/CRITICAL в собранном gateway. gosec выявил 27 HIGH-находок. Это
снимок состояния баз уязвимостей, числа будут меняться. Исправление и
обновление зависимостей — отдельная работа; отключать gates для получения
зелёного статуса нельзя.

Проверки идут последовательно и останавливаются на первой ошибке. Это помогает
понять причину красной сборки, но остальные сканеры в этом запуске не выполнятся.
После исправления перезапустите job. Подавлять находки следует только адресно,
с обоснованием в исходниках/настройках инструмента; общего baseline нет.

## Nexus: публикация после проверок

Nexus Repository Community Edition не нужен для обычного `Build Now`. При
первом запуске выберите бесплатный Community tier. Чтобы публиковать образы, в Nexus
создайте репозиторий **Docker (hosted)** с HTTP connector на порту `8082`,
включите **Docker Bearer Token Realm** в **Security → Realms** и создайте
пользователя с правом записи в этот репозиторий. В Jenkins добавьте credential
типа **Username with password** с ID `nexus-docker`. Затем запустите **Build
with Parameters** и установите `PUBLISH_IMAGES=true` только для проверенного
коммита `main`. Stage публикации расположен после всех quality gates.

Теги имеют вид `nexus:8082/volha/gateway:<номер-сборки>` и аналогично для
`auth-service`, `product-service`, `migrator`, `dumper`. На хосте тот же порт
открыт как `localhost:8087`; при pull вне Compose нужен соответствующий адрес
registry. В локальном dind явно разрешён HTTP registry `nexus:8082` — это
учебная конфигурация, для общего сервера настройте TLS. Учётные данные
передаются через Jenkins Credentials и `docker login --password-stdin` во
временный Docker config, который удаляется после публикации.

## GitHub Actions и дальнейший этап

Три workflow в `.github/workflows/` работают при push в `main`; backend workflow
также запускается для pull request. Они дают быстрый сигнал на GitHub и не
требуют работающего локального Jenkins. Jenkins запускает более тяжёлый набор
сканеров и SonarQube. GitHub branch protection может требовать успешные
workflow как условие merge; локальный Jenkins без публичного статуса GitHub
сам по себе merge не блокирует.

На следующем этапе CDL добавим отдельный pipeline с Testcontainers или
PostgreSQL/Redis, миграциями, API/E2E и позднее развёртыванием. Этот CI
сознательно не запускает `make test-integration`.

## Если что-то не стартует

- `docker compose config` требует `SONAR_DB_PASSWORD`: заполните
  `ci/jenkins/.env` и передайте `--env-file` как в командах выше.
- Jenkins пишет `Cannot connect to Docker daemon`: проверьте сервис `docker`
  через `docker compose -f ci/jenkins/compose.yaml --env-file ci/jenkins/.env logs docker` и общий том
  сертификатов. Jenkins может стартовать раньше выдачи сертификатов; повторите
  сборку после готовности dind.
- SonarQube долго запускается: смотрите `logs sonarqube sonar-db`, проверьте
  выделенную память Docker и параметры хоста для Elasticsearch. На Linux обычно
  нужен `vm.max_map_count` не ниже 524288; для Docker Desktop настройка зависит
  от используемой виртуальной машины.
- `sonar-token` или `nexus-docker` не найден: создайте credentials с точными ID.
- Сканер не видит SonarQube: проверьте `sonarqube` в Compose и адрес
  `http://sonarqube:9000` внутри CI-сети; `localhost:9000` подходит только
  браузеру на хосте.
- Красный quality gate: откройте отчёт в указанном выше месте; проверьте
  найденный файл, пакет и правило, затем исправьте проблему и повторите сборку.
