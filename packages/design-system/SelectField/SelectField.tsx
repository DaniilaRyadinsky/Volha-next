'use client'

import { useId, type ReactNode, type SelectHTMLAttributes } from 'react'
import { ChevronDown } from 'lucide-react'
import clsx from 'clsx'
import { FormField, getFieldDescriptionId } from '../Field/Field'
import styles from './SelectField.module.scss'

export type SelectFieldProps = SelectHTMLAttributes<HTMLSelectElement> & {
  label?: ReactNode
  hint?: ReactNode
  error?: ReactNode
  fieldClassName?: string
}

export function SelectField({
  'aria-describedby': ariaDescribedBy,
  children,
  className,
  disabled,
  error,
  fieldClassName,
  hint,
  id,
  label,
  required,
  ...props
}: SelectFieldProps) {
  const generatedId = useId()
  const controlId = id ?? generatedId
  const descriptionId = getFieldDescriptionId(controlId, error, hint)
  const describedBy = [ariaDescribedBy, descriptionId].filter(Boolean).join(' ') || undefined

  return (
    <FormField
      className={fieldClassName}
      disabled={disabled}
      error={error}
      hint={hint}
      id={controlId}
      label={label}
      required={required}
    >
      <span className={styles.shell}>
        <select
          aria-describedby={describedBy}
          aria-invalid={error ? true : undefined}
          className={clsx(styles.select, Boolean(error) && styles.invalid, className)}
          disabled={disabled}
          id={controlId}
          required={required}
          {...props}
        >
          {children}
        </select>
        <ChevronDown aria-hidden="true" className={styles.icon} />
      </span>
    </FormField>
  )
}
