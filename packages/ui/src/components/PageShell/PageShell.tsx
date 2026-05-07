import type { ComponentPropsWithoutRef, ElementType } from "react";
import styles from "./PageShell.module.css";

export type PageShellVariant = "grid" | "constellation";

interface PageShellOwnProps {
	as?: ElementType;
	variant?: PageShellVariant;
	center?: boolean;
	className?: string;
}

export type PageShellProps<T extends ElementType = "main"> = PageShellOwnProps & Omit<ComponentPropsWithoutRef<T>, keyof PageShellOwnProps>;

export function PageShell<T extends ElementType = "main">({ as, variant = "grid", center = false, className, ...props }: PageShellProps<T>) {
	const Comp = as || "main";
	const classes = ["ns-grid-shell", variant === "constellation" && "ns-grid-shell--constellation", center && styles.center, className].filter(Boolean).join(" ");

	return <Comp className={classes} {...props} />;
}
