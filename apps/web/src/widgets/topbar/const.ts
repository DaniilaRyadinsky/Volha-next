export type Page = "home" | "catalog" | "product" | "about" | "diy" | "contacts" | "cart";

export const NAV = [
  { label: "Каталог", page: "catalog" as Page },
  { label: "DIY и идеи", page: "diy" as Page },
  { label: "О нас", page: "about" as Page },
  { label: "Контакты", page: "contacts" as Page },
];