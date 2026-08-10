import type { HTMLAttributes, Key, ReactNode } from 'react'
import clsx from 'clsx'
import styles from './InstructionSteps.module.scss'

export type InstructionStep = {
  id?: Key
  label: ReactNode
}

export type InstructionStepsProps = HTMLAttributes<HTMLDivElement> & {
  label?: ReactNode
  steps: readonly InstructionStep[]
}

export function InstructionSteps({
  className,
  label = 'Пошаговая инструкция',
  steps,
  ...props
}: InstructionStepsProps) {
  return (
    <div className={clsx(styles.root, className)} {...props}>
      {label && <p className={styles.label}>{label}</p>}
      <ol aria-label={typeof label === 'string' ? label : undefined} className={styles.list}>
        {steps.map((step, index) => (
          <li className={styles.step} key={step.id ?? index}>
            <span aria-hidden="true" className={styles.number}>{index + 1}</span>
            <span>{step.label}</span>
          </li>
        ))}
      </ol>
    </div>
  )
}
