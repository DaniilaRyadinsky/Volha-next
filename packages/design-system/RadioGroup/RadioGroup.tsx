'use client'

import { useId, type ReactNode } from 'react'
import clsx from 'clsx'
import styles from './RadioGroup.module.scss'

export type RadioOption = {
  value: string
  label: ReactNode
  description?: ReactNode
  disabled?: boolean
}

export type RadioGroupProps = {
  className?: string
  defaultValue?: string
  disabled?: boolean
  legend?: ReactNode
  name?: string
  onValueChange?: (value: string) => void
  options: RadioOption[]
  orientation?: 'horizontal' | 'vertical'
  value?: string
}

export function RadioGroup({
  className,
  defaultValue,
  disabled = false,
  legend,
  name,
  onValueChange,
  options,
  orientation = 'vertical',
  value,
}: RadioGroupProps) {
  const generatedName = useId()
  const groupName = name ?? generatedName

  return (
    <fieldset className={clsx(styles.group, styles[orientation], className)} disabled={disabled}>
      {legend && <legend className={styles.legend}>{legend}</legend>}
      <div className={styles.options}>
        {options.map((option) => (
          <label className={clsx(styles.option, option.disabled && styles.disabled)} key={option.value}>
            <input
              checked={value !== undefined ? value === option.value : undefined}
              defaultChecked={value === undefined ? defaultValue === option.value : undefined}
              disabled={option.disabled}
              name={groupName}
              onChange={() => onValueChange?.(option.value)}
              type="radio"
              value={option.value}
            />
            <span aria-hidden="true" className={styles.control}><span /></span>
            <span className={styles.copy}>
              <span className={styles.label}>{option.label}</span>
              {option.description && <span className={styles.description}>{option.description}</span>}
            </span>
          </label>
        ))}
      </div>
    </fieldset>
  )
}
