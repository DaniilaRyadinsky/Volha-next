'use client'

import { useState, type CSSProperties, type HTMLAttributes, type Key } from 'react'
import clsx from 'clsx'
import styles from './ThumbnailGallery.module.scss'

export type ThumbnailGalleryImage = {
  id?: Key
  src: string
  thumbnailSrc?: string
  alt: string
}

export type ThumbnailGalleryProps = Omit<HTMLAttributes<HTMLDivElement>, 'onChange'> & {
  images: readonly ThumbnailGalleryImage[]
  activeIndex?: number
  defaultActiveIndex?: number
  onActiveChange?: (index: number) => void
  aspectRatio?: string
  ariaLabel?: string
}

export function ThumbnailGallery({
  activeIndex,
  ariaLabel = 'Галерея проекта',
  aspectRatio = '4 / 3',
  className,
  defaultActiveIndex = 0,
  images,
  onActiveChange,
  style,
  ...props
}: ThumbnailGalleryProps) {
  const [internalIndex, setInternalIndex] = useState(defaultActiveIndex)
  const requestedIndex = activeIndex ?? internalIndex
  const selectedIndex = Math.min(Math.max(requestedIndex, 0), Math.max(images.length - 1, 0))
  const selected = images[selectedIndex]

  const galleryStyle = {
    ...style,
    '--volha-gallery-aspect': aspectRatio,
    '--volha-gallery-columns': Math.max(images.length, 1),
  } as CSSProperties

  function select(index: number) {
    if (activeIndex === undefined) setInternalIndex(index)
    onActiveChange?.(index)
  }

  return (
    <div
      aria-label={ariaLabel}
      className={clsx(styles.gallery, className)}
      role="region"
      style={galleryStyle}
      {...props}
    >
      {selected && (
        <div className={styles.main}>
          <img alt={selected.alt} src={selected.src} />
        </div>
      )}
      {images.length > 1 && (
        <div className={styles.thumbnails}>
          {images.map((image, index) => (
            <button
              aria-label={`Показать изображение ${index + 1}: ${image.alt}`}
              aria-pressed={index === selectedIndex}
              className={clsx(styles.thumbnail, index === selectedIndex && styles.active)}
              key={image.id ?? image.src}
              onClick={() => select(index)}
              type="button"
            >
              <img alt="" loading="lazy" src={image.thumbnailSrc ?? image.src} />
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
