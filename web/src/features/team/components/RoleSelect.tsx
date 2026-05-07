import { useState } from "react";
import { classNames } from "../../../shared/utils/classNames";
import styles from "./RoleSelect.module.css";

interface RoleSelectProps {
	role: string;
	name: string;
}

const roleOptions = [
	{ value: "owner", label: "Owner" },
	{ value: "admin", label: "Admin" },
	{ value: "member", label: "Member" },
	{ value: "viewer", label: "Viewer" }
];

export function RoleSelect({ role, name }: RoleSelectProps) {
	const [selectedRole, setSelectedRole] = useState(role.toLowerCase());
	const roleClass = styles[selectedRole as keyof typeof styles] || styles.member;

	return (
		<span className={classNames(styles.frame, roleClass)}>
			<select className={styles.select} value={selectedRole} aria-label={`Change role for ${name}`} onChange={event => setSelectedRole(event.currentTarget.value)}>
				{roleOptions.map(option => (
					<option key={option.value} value={option.value}>
						{option.label}
					</option>
				))}
			</select>
		</span>
	);
}
