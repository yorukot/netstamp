import type { ComponentPropsWithoutRef, ElementType } from "react";
import styles from "./Surface.module.css";

export type SurfaceTone = "glass" | "matte" | "deep" | "flat" | "accent" | "danger";
export type SurfaceCut = "xs" | "sm" | "md" | "lg";
export type SurfacePadding = "none" | "sm" | "md" | "lg";

interface SurfaceOwnProps {
	as?: ElementType;
	tone?: SurfaceTone;
	cut?: SurfaceCut;
	padding?: SurfacePadding;
	className?: string;
}

export type SurfaceProps<T extends ElementType = "div"> = SurfaceOwnProps & Omit<ComponentPropsWithoutRef<T>, keyof SurfaceOwnProps>;

const cutClasses: Record<SurfaceCut, string> = {
	xs: styles.cutXs,
	sm: styles.cutSm,
	md: styles.cutMd,
	lg: styles.cutLg
};

const paddingClasses: Record<SurfacePadding, string> = {
	none: styles.paddingNone,
	sm: styles.paddingSm,
	md: styles.paddingMd,
	lg: styles.paddingLg
};

export function Surface<T extends ElementType = "div">({ as, tone = "glass", cut = "md", padding = "md", className, ...props }: SurfaceProps<T>) {
	const Comp = as || "div";
	const classes = ["ns-cut-frame", styles.surface, styles[tone], cutClasses[cut], paddingClasses[padding], className].filter(Boolean).join(" ");

	return <Comp className={classes} {...props} />;
}
