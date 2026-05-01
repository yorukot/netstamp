import type { ComponentPropsWithoutRef } from 'react'
import styles from './Badge.module.css'

export type BadgeVariant = 'default' | 'secondary' | 'outline' | 'success' | 'danger'

export interface BadgeProps extends ComponentPropsWithoutRef<'span'> {
  variant?: BadgeVariant
}

export function Badge({ variant = 'default', className, ...props }: BadgeProps) {
  const classes = [styles.badge, styles[variant], className].filter(Boolean).join(' ')
  return <span className={classes} {...props} />
}
