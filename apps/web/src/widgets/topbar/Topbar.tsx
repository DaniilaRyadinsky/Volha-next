'use client'

import React, { useEffect, useState } from 'react'
import styles from './Topbar.module.scss'
import { VolhaLogo } from '@volha/design-system'
import { ShoppingBag } from 'lucide-react'
import { NAV, Page } from './const'

export const Topbar = ({
  current,
  cartCount,
}: {
  current: Page;
  cartCount: number
}) => {
  const [open, setOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  // const { user, logout } = useAuth();
  const light = !scrolled && current === "home";

  useEffect(() => {
    const h = () => setScrolled(window.scrollY > 30);
    window.addEventListener("scroll", h);
    return () => window.removeEventListener("scroll", h);
  }, []);

  return (
    <div className={styles.patternNav}>
      <VolhaLogo light />
      <nav>
        {NAV.map((item) => (
          <a
            key={item.page}
            href={item.page === "home" ? "/" : `/${item.page}`}
            className={current === item.page ? styles.active : ''}
          >
            {item.label}
          </a>
        ))}
      </nav>
      <ShoppingBag />
    </div>
  )
}
