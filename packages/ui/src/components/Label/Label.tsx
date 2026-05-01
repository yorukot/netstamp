import * as LabelPrimitive from '@radix-ui/react-label'
import type { ComponentPropsWithoutRef } from 'react'
import styles from './Label.module.css'

export type LabelProps = ComponentPropsWithoutRef<typeof LabelPrimitive.Root>

export function Label({ className, ...props }: LabelProps) {
  const classes = [styles.label, className].filter(Boolean).join(' ')
  return <LabelPrimitive.Root className={classes} {...props} />
}
