'use client'

import { useEffect, useId, useRef, type InputHTMLAttributes, type ReactNode } from 'react'
import { Check, Minus } from 'lucide-react'
import clsx from 'clsx'
import styles from './Checkbox.module.scss'

export type CheckboxProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'type' | 'ref'> & {
  children: ReactNode
  description?: ReactNode
  indeterminate?: boolean
}

export function Checkbox({
  children,
  className,
  description,
  disabled,
  id,
  indeterminate = false,
  ...props
}: CheckboxProps) {
  const generatedId = useId()
  const controlId = id ?? generatedId
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (inputRef.current) inputRef.current.indeterminate = indeterminate
  }, [indeterminate])

  return (
    <label className={clsx(styles.checkbox, disabled && styles.disabled, className)} htmlFor={controlId}>
      <input disabled={disabled} id={controlId} ref={inputRef} type="checkbox" {...props} />
      <span aria-hidden="true" className={styles.control}>
        {indeterminate ? <Minus /> : <Check />}
      </span>
      <span className={styles.copy}>
        <span className={styles.label}>{children}</span>
        {description && <span className={styles.description}>{description}</span>}
      </span>
    </label>
  )
}
