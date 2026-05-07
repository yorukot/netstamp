import type { ComponentPropsWithoutRef } from "react";
import { classNames } from "../utils/classNames";
import styles from "./PageStack.module.css";

type PageStackProps = ComponentPropsWithoutRef<"section">;

export function PageStack({ className, ...props }: PageStackProps) {
	return <section className={classNames(styles.root, className)} {...props} />;
}
