import * as Label from "@radix-ui/react-label";
import type { ComponentPropsWithoutRef, ReactNode } from "react";
import { useId } from "react";
import styles from "./Field.module.css";

export interface FieldShellProps {
	id: string;
	label: ReactNode;
	helper?: ReactNode;
	error?: ReactNode;
	children: ReactNode;
}

function FieldShell({ id, label, helper, error, children }: FieldShellProps) {
	return (
		<div className={styles.field}>
			<Label.Root className={styles.label} htmlFor={id}>
				{label}
			</Label.Root>
			{children}
			{error ? <span className={styles.error}>{error}</span> : null}
			{helper && !error ? <span className={styles.helper}>{helper}</span> : null}
		</div>
	);
}

export interface TextFieldProps extends ComponentPropsWithoutRef<"input"> {
	label: ReactNode;
	helper?: ReactNode;
	error?: ReactNode;
}

export function TextField({ label, helper, error, className, ...props }: TextFieldProps) {
	const generatedId = useId();
	const id = props.id || generatedId;
	const classes = [styles.control, className].filter(Boolean).join(" ");

	return (
		<FieldShell id={id} label={label} helper={helper} error={error}>
			<span className={["ns-cut-frame", styles.controlFrame].join(" ")} data-invalid={Boolean(error)}>
				<input id={id} className={classes} aria-invalid={Boolean(error)} {...props} />
			</span>
		</FieldShell>
	);
}

export interface TextAreaFieldProps extends ComponentPropsWithoutRef<"textarea"> {
	label: ReactNode;
	helper?: ReactNode;
	error?: ReactNode;
}

export function TextAreaField({ label, helper, error, className, ...props }: TextAreaFieldProps) {
	const generatedId = useId();
	const id = props.id || generatedId;
	const classes = [styles.control, styles.area, className].filter(Boolean).join(" ");

	return (
		<FieldShell id={id} label={label} helper={helper} error={error}>
			<span className={["ns-cut-frame", styles.controlFrame].join(" ")} data-invalid={Boolean(error)}>
				<textarea id={id} className={classes} aria-invalid={Boolean(error)} {...props} />
			</span>
		</FieldShell>
	);
}

export interface SelectFieldProps extends ComponentPropsWithoutRef<"select"> {
	label: ReactNode;
	helper?: ReactNode;
	error?: ReactNode;
	options: Array<{ value: string; label: string }>;
}

export function SelectField({ label, helper, error, options, className, ...props }: SelectFieldProps) {
	const generatedId = useId();
	const id = props.id || generatedId;
	const classes = [styles.control, styles.select, className].filter(Boolean).join(" ");
	const frameClasses = ["ns-cut-frame", styles.controlFrame, styles.selectFrame].join(" ");

	return (
		<FieldShell id={id} label={label} helper={helper} error={error}>
			<span className={frameClasses} data-invalid={Boolean(error)}>
				<select id={id} className={classes} aria-invalid={Boolean(error)} {...props}>
					{options.map(option => (
						<option key={option.value} value={option.value}>
							{option.label}
						</option>
					))}
				</select>
			</span>
		</FieldShell>
	);
}
