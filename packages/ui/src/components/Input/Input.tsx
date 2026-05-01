import type { ComponentPropsWithoutRef } from 'react'
import styles from './Input.module.css'

export interface InputProps extends ComponentPropsWithoutRef<'input'> {
  invalid?: boolean
}

export function Input({ invalid, className, ...props }: InputProps) {
  const classes = [styles.input, className].filter(Boolean).join(' ')
  return (
    <input
      className={classes}
      aria-invalid={invalid || undefined}
      {...props}
    />
  )
}
