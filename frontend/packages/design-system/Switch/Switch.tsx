'use client'

import { useId, type InputHTMLAttributes, type ReactNode } from 'react'
import clsx from 'clsx'
import styles from './Switch.module.scss'

export type SwitchProps = Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> & {
  label?: ReactNode
  description?: ReactNode
}

export function Switch({ className, description, disabled, id, label, ...props }: SwitchProps) {
  const generatedId = useId()
  const controlId = id ?? generatedId

  return (
    <label className={clsx(styles.switch, disabled && styles.disabled, className)} htmlFor={controlId}>
      <input disabled={disabled} id={controlId} role="switch" type="checkbox" {...props} />
      <span aria-hidden="true" className={styles.track}><span className={styles.thumb} /></span>
      {(label || description) && (
        <span className={styles.copy}>
          {label && <span className={styles.label}>{label}</span>}
          {description && <span className={styles.description}>{description}</span>}
        </span>
      )}
    </label>
  )
}
