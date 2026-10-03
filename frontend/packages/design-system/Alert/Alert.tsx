'use client'

import type { HTMLAttributes, ReactNode } from 'react'
import { CircleCheck, CircleX, Info, TriangleAlert, X } from 'lucide-react'
import clsx from 'clsx'
import { IconButton } from '../IconButton/IconButton'
import styles from './Alert.module.scss'

export type AlertTone = 'info' | 'success' | 'warning' | 'danger'

export type AlertProps = Omit<HTMLAttributes<HTMLDivElement>, 'title'> & {
  title?: ReactNode
  tone?: AlertTone
  icon?: ReactNode
  action?: ReactNode
  onDismiss?: () => void
  dismissLabel?: string
}

const toneIcons = {
  info: <Info />,
  success: <CircleCheck />,
  warning: <TriangleAlert />,
  danger: <CircleX />,
}

export function Alert({
  action,
  children,
  className,
  dismissLabel = 'Dismiss',
  icon,
  onDismiss,
  role,
  title,
  tone = 'info',
  ...props
}: AlertProps) {
  return (
    <div className={clsx(styles.alert, styles[tone], className)} role={role ?? (tone === 'danger' ? 'alert' : 'status')} {...props}>
      <span aria-hidden="true" className={styles.icon}>{icon ?? toneIcons[tone]}</span>
      <div className={styles.content}>
        {title && <strong>{title}</strong>}
        {children && <div className={styles.message}>{children}</div>}
      </div>
      {action && <div className={styles.action}>{action}</div>}
      {onDismiss && (
        <IconButton className={styles.dismiss} label={dismissLabel} onClick={onDismiss} size="sm" variant="ghost">
          <X />
        </IconButton>
      )}
    </div>
  )
}
