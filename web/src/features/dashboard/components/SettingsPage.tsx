import { Button, Panel, TextField } from "@netstamp/ui";
import type { FormEvent } from "react";
import { ScreenHeader } from "../../../shared/components/ScreenHeader";
import { currentUser } from "../../../shared/utils/mockData";
import styles from "./ProductPages.module.css";

function handleSettingsSubmit(event: FormEvent<HTMLFormElement>) {
	event.preventDefault();
}

export function SettingsPage() {
	return (
		<section className={styles.screen}>
			<ScreenHeader eyebrow="User settings" title="Account" copy="Set your username, rotate the login email, and change the password used for controller access." />

			<div className={styles.settingsGrid}>
				<Panel tone="glass" eyebrow="Identity" title="Set username">
					<form id="username-settings" className={styles.settingsForm} onSubmit={handleSettingsSubmit}>
						<TextField label="Display name" name="name" defaultValue={currentUser.name} />
						<TextField label="Username" name="username" defaultValue={currentUser.username} helper="Used in audit events and probe ownership trails." />
						<div className={styles.actionRow}>
							<Button type="submit">Save username</Button>
						</div>
					</form>
				</Panel>

				<Panel tone="deep" eyebrow="Profile image" title="Gravatar signal preview">
					<div className={styles.profilePreview}>
						<span className={styles.profileFrame} aria-hidden="true">
							<img src={currentUser.gravatarUrl} alt="" referrerPolicy="no-referrer" />
						</span>
						<div>
							<h3>{currentUser.name}</h3>
							<p>{currentUser.email}</p>
						</div>
					</div>
					<p className={styles.bodyCopy}>The avatar is pulled using your email from Gravatar.</p>
				</Panel>
			</div>

			<div className={styles.settingsGrid}>
				<Panel tone="glass" eyebrow="Email" title="Change email">
					<form className={styles.settingsForm} onSubmit={handleSettingsSubmit}>
						<TextField label="Current email" name="current-email" type="email" defaultValue={currentUser.email} />
						<TextField label="New email" name="new-email" type="email" placeholder="operator@example.com" />
						<TextField label="Confirm password" name="email-password" type="password" autoComplete="current-password" />
						<div className={styles.actionRow}>
							<Button type="submit">Update email</Button>
						</div>
					</form>
				</Panel>

				<Panel tone="glass" eyebrow="Security" title="Change password">
					<form className={styles.settingsForm} onSubmit={handleSettingsSubmit}>
						<TextField label="Current password" name="current-password" type="password" autoComplete="current-password" />
						<TextField label="New password" name="new-password" type="password" autoComplete="new-password" />
						<TextField label="Confirm new password" name="confirm-password" type="password" autoComplete="new-password" helper="Use at least 12 characters for production accounts." />
						<div className={styles.actionRow}>
							<Button type="submit">Change password</Button>
						</div>
					</form>
				</Panel>
			</div>
		</section>
	);
}
