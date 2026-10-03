import type { AnchorHTMLAttributes, ReactNode } from 'react'
import { ArrowLeft } from 'lucide-react'
import clsx from 'clsx'
import styles from './BackLink.module.scss'

export type BackLinkProps = AnchorHTMLAttributes<HTMLAnchorElement> & {
  icon?: ReactNode
}

export function BackLink({ children, className, icon, ...props }: BackLinkProps) {
  return (
    <a className={clsx(styles.link, className)} {...props}>
      <span aria-hidden="true" className={styles.icon}>
        {icon ?? <ArrowLeft />}
      </span>
      <span>{children}</span>
    </a>
  )
}
