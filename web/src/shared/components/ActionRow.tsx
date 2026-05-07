import { classNames } from "@/shared/utils/classNames";
import type { ComponentPropsWithoutRef } from "react";
import styles from "./ActionRow.module.css";

type ActionRowProps = ComponentPropsWithoutRef<"div">;

export function ActionRow({ className, ...props }: ActionRowProps) {
	return <div className={classNames(styles.root, className)} {...props} />;
}
