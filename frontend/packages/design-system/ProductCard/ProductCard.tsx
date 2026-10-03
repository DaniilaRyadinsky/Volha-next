'use client'

import type { HTMLAttributes, ReactNode } from 'react'
import { ShoppingBag } from 'lucide-react'
import { Badge, type BadgeTone } from '../Badge/Badge'
import clsx from 'clsx'
import { IconButton } from '../IconButton/IconButton'
import styles from './ProductCard.module.scss'

export type ProductCardVariant = 'default' | 'project'

export type ProductCardProps = Omit<HTMLAttributes<HTMLElement>, 'title'> & {
  image: string
  imageAlt?: string
  category: string
  title: string
  price: string
  badge?: ReactNode
  badgeTone?: BadgeTone
  onAdd?: () => void
  addLabel?: string
  actionDisabled?: boolean
  variant?: ProductCardVariant
}

export function ProductCard({
  actionDisabled = false,
  addLabel,
  badge,
  badgeTone = 'accent',
  category,
  className,
  image,
  imageAlt,
  onAdd,
  price,
  title,
  variant = 'default',
  ...props
}: ProductCardProps) {
  return (
    <article className={clsx(styles.card, variant === 'project' && styles.project, className)} {...props}>
      <div className={styles.media}>
        <img alt={imageAlt ?? title} loading="lazy" src={image} />
        {badge && <Badge className={styles.badge} size="sm" tone={badgeTone}>{badge}</Badge>}
      </div>
      <div className={styles.body}>
        <p className={styles.category}>{category}</p>
        <h3>{title}</h3>
        <div className={styles.meta}>
          <span>{price}</span>
          {onAdd && (
            <IconButton
              className={styles.add}
              disabled={actionDisabled}
              label={addLabel ?? `Добавить ${title} в корзину`}
              onClick={onAdd}
              size="sm"
              variant="ghost"
            >
              <ShoppingBag />
            </IconButton>
          )}
        </div>
      </div>
    </article>
  )
}
