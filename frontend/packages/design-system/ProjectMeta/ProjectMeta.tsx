import type { HTMLAttributes, ReactNode } from 'react'
import { Clock3 } from 'lucide-react'
import clsx from 'clsx'
import { Badge, type BadgeTone } from '../Badge/Badge'
import styles from './ProjectMeta.module.scss'

export type ProjectMetaProps = HTMLAttributes<HTMLDivElement> & {
  difficulty: ReactNode
  difficultyTone?: BadgeTone
  category: ReactNode
  duration: ReactNode
  durationIcon?: ReactNode
}

export function ProjectMeta({
  category,
  className,
  difficulty,
  difficultyTone = 'medium',
  duration,
  durationIcon,
  ...props
}: ProjectMetaProps) {
  return (
    <div className={clsx(styles.meta, className)} {...props}>
      <Badge size="sm" tone={difficultyTone}>{difficulty}</Badge>
      <span className={styles.item}>{category}</span>
      <span className={styles.item}>
        <span aria-hidden="true" className={styles.icon}>{durationIcon ?? <Clock3 />}</span>
        {duration}
      </span>
    </div>
  )
}
