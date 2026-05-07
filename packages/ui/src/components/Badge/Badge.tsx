import type { ComponentPropsWithoutRef } from "react";
import styles from "./Badge.module.css";

export type BadgeTone = "neutral" | "accent" | "success" | "warning" | "critical" | "muted";

export interface BadgeProps extends ComponentPropsWithoutRef<"span"> {
	tone?: BadgeTone;
	dot?: boolean;
}

export function Badge({ tone = "neutral", dot = true, className, children, ...props }: BadgeProps) {
	const classes = ["ns-cut-frame", styles.badge, styles[tone], className].filter(Boolean).join(" ");

	return (
		<span className={classes} {...props}>
			{dot ? <span className={styles.dot} aria-hidden="true" /> : null}
			{children}
		</span>
	);
}
