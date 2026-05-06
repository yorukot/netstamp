import { Badge, Panel } from '@netstamp/ui'
import { systemStates } from '../utils/mockData'
import styles from './SystemStateGrid.module.css'

export function SystemStateGrid() {
  return (
    <div className={styles.grid}>
      {systemStates.map(([title, copy], index) => (
        <Panel key={title} tone={index % 3 === 0 ? 'deep' : 'glass'}>
          <div className={styles.state}>
            <Badge tone={index < 2 ? 'accent' : index > 4 ? 'critical' : 'neutral'}>
              {String(index + 1).padStart(2, '0')}
            </Badge>
            <h3>{title}</h3>
            <p>{copy}</p>
            {index < 2 ? <span className={styles.loader} /> : null}
          </div>
        </Panel>
      ))}
    </div>
  )
}
