import type { ReactNode } from 'react'
import styles from './DataTable.module.css'

export interface DataColumn<Row extends object = Record<string, unknown>> {
  key: string
  label: ReactNode
  render?: (row: Row, index: number) => ReactNode
}

export interface DataTableProps<Row extends object = Record<string, unknown>> {
  columns: DataColumn<Row>[]
  rows: Row[]
  getRowKey?: (row: Row, index: number) => string
  onRowClick?: (row: Row) => void
  selectedKey?: string
  emptyLabel?: ReactNode
}

export function DataTable<Row extends object>({
  columns,
  rows,
  getRowKey,
  onRowClick,
  selectedKey,
  emptyLabel = 'No results',
}: DataTableProps<Row>) {
  return (
    <div className={styles.wrap}>
      <table className={styles.table}>
        <thead>
          <tr>
            {columns.map((column) => (
              <th key={column.key}>{column.label}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {rows.length ? (
            rows.map((row, index) => {
              const rowKey = getRowKey ? getRowKey(row, index) : String(index)
              const selected = selectedKey === rowKey

              return (
                <tr
                  key={rowKey}
                  className={selected ? styles.selected : undefined}
                  onClick={onRowClick ? () => onRowClick(row) : undefined}
                >
                  {columns.map((column) => (
                    <td key={column.key}>
                      {column.render ? column.render(row, index) : String((row as Record<string, unknown>)[column.key] ?? '')}
                    </td>
                  ))}
                </tr>
              )
            })
          ) : (
            <tr>
              <td className={styles.empty} colSpan={columns.length}>
                {emptyLabel}
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  )
}
