import type { ReactNode } from "react";
import styles from "./ScreenHeader.module.css";

interface ScreenHeaderProps {
	eyebrow?: ReactNode;
	title: ReactNode;
	copy?: ReactNode;
	actions?: ReactNode;
}

export function ScreenHeader({ eyebrow, title, copy, actions }: ScreenHeaderProps) {
	return (
		<header className={styles.header}>
			<div>
				{eyebrow ? <span className={styles.eyebrow}>{eyebrow}</span> : null}
				<h1>{title}</h1>
				{copy ? <p>{copy}</p> : null}
			</div>
			{actions ? <div className={styles.actions}>{actions}</div> : null}
		</header>
	);
}
