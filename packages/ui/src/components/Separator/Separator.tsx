import * as SeparatorPrimitive from '@radix-ui/react-separator'
import type { ComponentPropsWithoutRef } from 'react'
import styles from './Separator.module.css'

export type SeparatorProps = ComponentPropsWithoutRef<typeof SeparatorPrimitive.Root>

export function Separator({ className, ...props }: SeparatorProps) {
  const classes = [styles.separator, className].filter(Boolean).join(' ')
  return <SeparatorPrimitive.Root className={classes} {...props} />
}
