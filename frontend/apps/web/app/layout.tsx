import type { Metadata } from 'next'
import { Topbar } from '@/widgets/topbar'

export const metadata: Metadata = {
  title: 'Volha UI — дизайн-система',
  description: 'Компоненты и токены интернет-магазина гаражного оборудования Volha',
}

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="ru">
      <body>
        <Topbar current="home" cartCount={0} />
        {children}</body>
    </html>
  )
}
