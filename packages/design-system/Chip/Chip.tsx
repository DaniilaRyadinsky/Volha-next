'use client'

import type { ButtonHTMLAttributes, ReactNode } from 'react'
import { X } from 'lucide-react'
import clsx from 'clsx'
import styles from './Chip.module.scss'

export type ChipProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  active?: boolean
  leadingIcon?: ReactNode
  onRemove?: () => void
  removeLabel?: string
}

export function Chip({
  active = false,
  children,
  className,
  leadingIcon,
  onRemove,
  removeLabel = 'Remove',
  type = 'button',
  ...props
}: ChipProps) {
  return (
    <span className={clsx(styles.shell, active && styles.active, className)}>
      <button aria-pressed={active} className={styles.action} type={type} {...props}>
        {leadingIcon && <span className={styles.icon}>{leadingIcon}</span>}
        <span>{children}</span>
      </button>
      {onRemove && (
        <button aria-label={removeLabel} className={styles.remove} onClick={onRemove} title={removeLabel} type="button">
          <X />
        </button>
      )}
    </span>
  )
}

export const FilterChip = Chip
