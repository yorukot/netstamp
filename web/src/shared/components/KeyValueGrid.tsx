import type { ReactNode } from "react";
import { classNames } from "../utils/classNames";
import styles from "./KeyValueGrid.module.css";

export interface KeyValueItem {
	key?: string;
	label: ReactNode;
	value: ReactNode;
}

interface KeyValueGridProps {
	className?: string;
	items: KeyValueItem[];
}

export function KeyValueGrid({ className, items }: KeyValueGridProps) {
	return (
		<div className={classNames(styles.grid, className)}>
			{items.map((item, index) => (
				<div className={styles.card} key={item.key ?? (typeof item.label === "string" ? item.label : index)}>
					<span className={styles.label}>{item.label}</span>
					<strong className={styles.value}>{item.value}</strong>
				</div>
			))}
		</div>
	);
}
