'use client'

import { useId, type ReactNode, type TextareaHTMLAttributes } from 'react'
import clsx from 'clsx'
import { FormField, getFieldDescriptionId } from '../Field/Field'
import styles from './TextareaField.module.scss'

export type TextareaFieldProps = TextareaHTMLAttributes<HTMLTextAreaElement> & {
  label?: ReactNode
  hint?: ReactNode
  error?: ReactNode
  fieldClassName?: string
}

export function TextareaField({
  'aria-describedby': ariaDescribedBy,
  className,
  disabled,
  error,
  fieldClassName,
  hint,
  id,
  label,
  required,
  ...props
}: TextareaFieldProps) {
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
      <textarea
        aria-describedby={describedBy}
        aria-invalid={error ? true : undefined}
          className={clsx(styles.textarea, Boolean(error) && styles.invalid, className)}
        disabled={disabled}
        id={controlId}
        required={required}
        {...props}
      />
    </FormField>
  )
}
