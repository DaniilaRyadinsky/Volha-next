'use client'

import type { ReactNode } from 'react'
import clsx from 'clsx'
import styles from './SegmentedControl.module.scss'

export type SegmentedOption<T extends string> = {
  value: T
  label: ReactNode
  icon?: ReactNode
  disabled?: boolean
}

export type SegmentedControlProps<T extends string> = {
  'aria-label': string
  className?: string
  fullWidth?: boolean
  onValueChange: (value: T) => void
  options: SegmentedOption<T>[]
  size?: 'sm' | 'md'
  value: T
}

export function SegmentedControl<T extends string>({
  'aria-label': ariaLabel,
  className,
  fullWidth = false,
  onValueChange,
  options,
  size = 'md',
  value,
}: SegmentedControlProps<T>) {
  return (
    <div aria-label={ariaLabel} className={clsx(styles.control, styles[size], fullWidth && styles.fullWidth, className)} role="group">
      {options.map((option) => (
        <button
          aria-pressed={value === option.value}
          className={clsx(styles.option, value === option.value && styles.active)}
          disabled={option.disabled}
          key={option.value}
          onClick={() => onValueChange(option.value)}
          type="button"
        >
          {option.icon && <span className={styles.icon}>{option.icon}</span>}
          <span>{option.label}</span>
        </button>
      ))}
    </div>
  )
}
