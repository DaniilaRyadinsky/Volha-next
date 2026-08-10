'use client'

import { useId, type InputHTMLAttributes, type ReactNode } from 'react'
import clsx from 'clsx'
import { FormField, getFieldDescriptionId } from '../Field/Field'
import fieldStyles from '../Field/Field.module.scss'
import styles from './InputField.module.scss'

export type InputFieldProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'size'> & {
  label?: ReactNode
  hint?: ReactNode
  error?: ReactNode
  leadingIcon?: ReactNode
  icon?: ReactNode
  trailingElement?: ReactNode
  fieldClassName?: string
}

export function InputField({
  'aria-describedby': ariaDescribedBy,
  className,
  disabled,
  error,
  fieldClassName,
  hint,
  id,
  icon,
  label,
  leadingIcon,
  required,
  trailingElement,
  ...props
}: InputFieldProps) {
  const generatedId = useId()
  const controlId = id ?? generatedId
  const startIcon = leadingIcon ?? icon
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
      <span className={fieldStyles.controlShell}>
        {startIcon && <span className={fieldStyles.leading}>{startIcon}</span>}
        <input
          aria-describedby={describedBy}
          aria-invalid={error ? true : undefined}
          className={clsx(styles.input, Boolean(startIcon) && styles.withLeading, Boolean(trailingElement) && styles.withTrailing, Boolean(error) && styles.invalid, className)}
          disabled={disabled}
          id={controlId}
          required={required}
          {...props}
        />
        {trailingElement && <span className={fieldStyles.trailing}>{trailingElement}</span>}
      </span>
    </FormField>
  )
}

export const Field = InputField
