'use client'

import { useEffect, useId, useRef, useState, type ReactNode, type RefObject } from 'react'
import { createPortal } from 'react-dom'
import { X } from 'lucide-react'
import clsx from 'clsx'
import { IconButton } from '../IconButton/IconButton'
import styles from './Modal.module.scss'

const focusableSelector = [
  'a[href]',
  'button:not([disabled])',
  'input:not([disabled])',
  'select:not([disabled])',
  'textarea:not([disabled])',
  '[tabindex]:not([tabindex="-1"])',
].join(',')

export type ModalProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  title?: ReactNode
  description?: ReactNode
  children: ReactNode
  footer?: ReactNode
  ariaLabel?: string
  className?: string
  size?: 'sm' | 'md' | 'lg' | 'full'
  closeLabel?: string
  closeOnBackdrop?: boolean
  closeOnEscape?: boolean
  showClose?: boolean
  initialFocusRef?: RefObject<HTMLElement | null>
}

export function Modal({
  ariaLabel,
  children,
  className,
  closeLabel = 'Close dialog',
  closeOnBackdrop = true,
  closeOnEscape = true,
  description,
  footer,
  initialFocusRef,
  onOpenChange,
  open,
  showClose = true,
  size = 'md',
  title,
}: ModalProps) {
  const [mounted, setMounted] = useState(false)
  const dialogRef = useRef<HTMLDivElement>(null)
  const titleId = useId()
  const descriptionId = useId()

  useEffect(() => setMounted(true), [])

  useEffect(() => {
    if (!open || !mounted) return

    const previousActiveElement = document.activeElement instanceof HTMLElement ? document.activeElement : null
    const previousOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'

    const dialog = dialogRef.current
    const preferredFocus = initialFocusRef?.current
    const firstFocusable = dialog?.querySelector<HTMLElement>(focusableSelector)
    window.requestAnimationFrame(() => (preferredFocus ?? firstFocusable ?? dialog)?.focus())

    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape' && closeOnEscape) {
        event.preventDefault()
        onOpenChange(false)
        return
      }

      if (event.key !== 'Tab' || !dialog) return
      const focusable = Array.from(dialog.querySelectorAll<HTMLElement>(focusableSelector))
        .filter((element) => !element.hasAttribute('disabled') && element.getAttribute('aria-hidden') !== 'true')

      if (focusable.length === 0) {
        event.preventDefault()
        dialog.focus()
        return
      }

      const first = focusable[0]
      const last = focusable[focusable.length - 1]
      if (event.shiftKey && document.activeElement === first) {
        event.preventDefault()
        last.focus()
      } else if (!event.shiftKey && document.activeElement === last) {
        event.preventDefault()
        first.focus()
      }
    }

    document.addEventListener('keydown', handleKeyDown)
    return () => {
      document.removeEventListener('keydown', handleKeyDown)
      document.body.style.overflow = previousOverflow
      previousActiveElement?.focus()
    }
  }, [closeOnEscape, initialFocusRef, mounted, onOpenChange, open])

  if (!mounted || !open) return null

  return createPortal(
    <div
      className={styles.overlay}
      onMouseDown={(event) => {
        if (closeOnBackdrop && event.target === event.currentTarget) onOpenChange(false)
      }}
    >
      <div
        aria-describedby={description ? descriptionId : undefined}
        aria-label={!title ? ariaLabel : undefined}
        aria-labelledby={title ? titleId : undefined}
        aria-modal="true"
        className={clsx(styles.dialog, styles[size], className)}
        ref={dialogRef}
        role="dialog"
        tabIndex={-1}
      >
        {(title || description || showClose) && (
          <header className={styles.header}>
            <div className={styles.heading}>
              {title && <h2 id={titleId}>{title}</h2>}
              {description && <p id={descriptionId}>{description}</p>}
            </div>
            {showClose && (
              <IconButton className={styles.close} label={closeLabel} onClick={() => onOpenChange(false)} size="sm" variant="ghost">
                <X />
              </IconButton>
            )}
          </header>
        )}
        <div className={styles.body}>{children}</div>
        {footer && <footer className={styles.footer}>{footer}</footer>}
      </div>
    </div>,
    document.body,
  )
}
