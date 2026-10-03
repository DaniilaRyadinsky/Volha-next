import type { ReactNode } from 'react'
import clsx from 'clsx'
import styles from './Field.module.scss'

export type FormFieldProps = {
  children: ReactNode
  className?: string
  disabled?: boolean
  error?: ReactNode
  hint?: ReactNode
  id: string
  label?: ReactNode
  required?: boolean
}

export function getFieldDescriptionId(id: string, error?: ReactNode, hint?: ReactNode) {
  if (error) return `${id}-error`
  if (hint) return `${id}-hint`
  return undefined
}

export function FormField({
  children,
  className,
  disabled = false,
  error,
  hint,
  id,
  label,
  required = false,
}: FormFieldProps) {
  return (
    <div className={clsx(styles.field, disabled && styles.disabled, className)}>
      {label && (
        <label className={styles.label} htmlFor={id}>
          <span>{label}</span>
          {required && <span aria-hidden="true" className={styles.required}>*</span>}
        </label>
      )}
      {children}
      {error ? (
        <span className={clsx(styles.message, styles.error)} id={`${id}-error`} role="alert">
          {error}
        </span>
      ) : hint ? (
        <span className={styles.message} id={`${id}-hint`}>
          {hint}
        </span>
      ) : null}
    </div>
  )
}
