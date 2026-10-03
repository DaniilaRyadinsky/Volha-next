'use client'

import type { ButtonHTMLAttributes, ReactNode } from 'react'
import clsx from 'clsx'
import { Spinner } from '../Spinner/Spinner'
import styles from './Button.module.scss'

export type ButtonVariant =
  | 'primary'
  | 'solid'
  | 'secondary'
  | 'quiet'
  | 'outline'
  | 'accent'
  | 'ghost'
  | 'light'
  | 'danger'

export type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant
  size?: 'sm' | 'md' | 'lg'
  fullWidth?: boolean
  loading?: boolean
  leadingIcon?: ReactNode
  trailingIcon?: ReactNode
  icon?: ReactNode
}

const variantClasses: Record<ButtonVariant, string> = {
  primary: styles.primary,
  solid: styles.primary,
  secondary: styles.secondary,
  quiet: styles.secondary,
  outline: styles.outline,
  accent: styles.accent,
  ghost: styles.ghost,
  light: styles.light,
  danger: styles.danger,
}

export function Button({
  children,
  className,
  disabled,
  fullWidth = false,
  icon,
  leadingIcon,
  loading = false,
  size = 'md',
  trailingIcon,
  type = 'button',
  variant = 'primary',
  ...props
}: ButtonProps) {
  const endIcon = trailingIcon ?? icon

  return (
    <button
      aria-busy={loading || undefined}
      className={clsx(
        styles.button,
        styles[size],
        variantClasses[variant],
        fullWidth && styles.fullWidth,
        loading && styles.loading,
        className,
      )}
      disabled={disabled || loading}
      type={type}
      {...props}
    >
      {loading ? <Spinner label="Loading" size="sm" /> : leadingIcon && <span className={styles.icon}>{leadingIcon}</span>}
      {children && <span className={styles.label}>{children}</span>}
      {!loading && endIcon && <span className={styles.icon}>{endIcon}</span>}
    </button>
  )
}
