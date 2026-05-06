import type { ReactNode } from 'react'
import { Badge, type BadgeTone } from '../Badge/Badge'
import styles from './MetricCard.module.css'

export interface MetricCardProps {
  label: ReactNode
  value: ReactNode
  detail?: ReactNode
  tone?: BadgeTone
  className?: string
}

export function MetricCard({ label, value, detail, tone = 'accent', className }: MetricCardProps) {
  const classes = [styles.card, className].filter(Boolean).join(' ')

  return (
    <article className={classes}>
      <span className={styles.label}>{label}</span>
      <strong>{value}</strong>
      {detail ? <Badge tone={tone}>{detail}</Badge> : null}
    </article>
  )
}
