import type { ComponentPropsWithoutRef } from 'react'
import styles from './Card.module.css'

export function Card({ className, ...props }: ComponentPropsWithoutRef<'div'>) {
  return <div className={[styles.card, className].filter(Boolean).join(' ')} {...props} />
}

export function CardHeader({ className, ...props }: ComponentPropsWithoutRef<'div'>) {
  return <div className={[styles.header, className].filter(Boolean).join(' ')} {...props} />
}

export function CardTitle({ className, ...props }: ComponentPropsWithoutRef<'h3'>) {
  return <h3 className={[styles.title, className].filter(Boolean).join(' ')} {...props} />
}

export function CardDescription({ className, ...props }: ComponentPropsWithoutRef<'p'>) {
  return <p className={[styles.description, className].filter(Boolean).join(' ')} {...props} />
}

export function CardContent({ className, ...props }: ComponentPropsWithoutRef<'div'>) {
  return <div className={[styles.content, className].filter(Boolean).join(' ')} {...props} />
}

export function CardFooter({ className, ...props }: ComponentPropsWithoutRef<'div'>) {
  return <div className={[styles.footer, className].filter(Boolean).join(' ')} {...props} />
}
