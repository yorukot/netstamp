import { Button } from "@netstamp/ui";
import { Link, NavLink, Outlet } from "react-router-dom";
import { pathForRoute } from "../routes/routePaths";
import { GlobalFooter } from "../shared/components/GlobalFooter";
import { currentUser, sidebarItems } from "../shared/utils/mockData";
import styles from "./AppShell.module.css";

export function AppShell() {
	return (
		<div className={styles.shell}>
			<aside className={styles.sidebar}>
				<Link className={styles.brand} to={pathForRoute("landing")}>
					<span className={styles.mark} aria-hidden="true" />
					<span>Netstamp</span>
				</Link>

				<label className={styles.teamSelect}>
					<span>team</span>
					<span className={styles.teamFrame}>
						<select defaultValue="vector-ix">
							<option value="vector-ix">Vector IX / prod</option>
							<option value="helio">Helio Validators</option>
							<option value="lab">Lab Network</option>
						</select>
					</span>
				</label>

				<nav className={styles.nav} aria-label="Primary app navigation">
					{sidebarItems.map(item => (
						<NavLink key={item.route} to={pathForRoute(item.route)} className={({ isActive }) => (isActive ? styles.active : undefined)}>
							{item.label}
						</NavLink>
					))}
				</nav>

				<div className={styles.userCard}>
					<div className={styles.userProfile}>
						<span className={styles.avatarFrame} aria-hidden="true">
							<img src={currentUser.gravatarUrl} alt="" referrerPolicy="no-referrer" />
						</span>
						<div className={styles.userMeta}>
							<strong>{currentUser.name}</strong>
							<span>{currentUser.role}</span>
						</div>
					</div>
					<div className={styles.userActions}>
						<Button variant="ghost" size="sm" asChild>
							<Link to={pathForRoute("landing")}>logout</Link>
						</Button>
						<Button variant="ghost" size="sm" asChild>
							<Link to={pathForRoute("settings")}>Settings</Link>
						</Button>
					</div>
				</div>
			</aside>

			<main className={styles.content}>
				<Outlet />
				<GlobalFooter variant="compact" />
			</main>
		</div>
	);
}
