import styles from "./FleetMatrix.module.css";

export function FleetMatrix({ total = 128, online = 100 }) {
	const cells = Array.from({ length: total }, (_, index) => index);

	return (
		<div className={styles.wrap} aria-label={`${online} of ${total} probes online`}>
			<div className={styles.header}>
				<span>fleet bitmap</span>
				<strong>
					{online}/{total}
				</strong>
			</div>
			<div className={styles.grid}>
				{cells.map(cell => (
					<span key={cell} className={cell < online ? styles.on : styles.off} />
				))}
			</div>
		</div>
	);
}
