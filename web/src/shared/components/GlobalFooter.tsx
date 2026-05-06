import { GithubLogoIcon } from "@phosphor-icons/react/dist/csr/GithubLogo";
import { NetworkIcon } from "@phosphor-icons/react/dist/csr/Network";
import { StarIcon } from "@phosphor-icons/react/dist/csr/Star";
import { Link } from "react-router-dom";
import { pathForRoute } from "../../routes/routePaths";
import styles from "./GlobalFooter.module.css";

const githubUrl = "https://github.com/yorukot/netstamp";

interface GlobalFooterProps {
	variant?: "full" | "compact";
	className?: string;
}

export function GlobalFooter({ variant = "full", className }: GlobalFooterProps) {
	const classes = [styles.footer, styles[variant], className].filter(Boolean).join(" ");

	return (
		<footer className={classes}>
			{variant === "full" ? (
				<div className={styles.footerGrid}>
					<div className={styles.footerBrand}>
						<span className={styles.brandMark} aria-hidden="true">
							<NetworkIcon size={18} weight="bold" />
						</span>
						<div>
							<strong>Netstamp</strong>
							<p>Open-source network measurement from probes you control.</p>
						</div>
					</div>

					<div className={styles.footerColumn}>
						<span>Product</span>
						<Link to={pathForRoute("dashboard")}>Console demo</Link>
						<Link to={pathForRoute("probes")}>Probe fleet</Link>
						<Link to={pathForRoute("components")}>Components</Link>
					</div>

					<div className={styles.footerColumn}>
						<span>Project</span>
						<a href={githubUrl} target="_blank" rel="noreferrer">
							GitHub source
						</a>
						<Link to={pathForRoute("register")}>Deploy a probe</Link>
						<Link to={pathForRoute("login")}>Operator login</Link>
					</div>
				</div>
			) : null}

			<div className={styles.footerBottom}>
				<span>
					Netstamp / Made by{" "}
					<a href="https://github.com/elvisdragonmao" target="_blank" rel="noreferrer">
						Elvis Mao
					</a>
					,{" "}
					<a href="https://github.com/yorukot" target="_blank" rel="noreferrer">
						Yorukot
					</a>
					, and{" "}
					<a href={githubUrl} target="_blank" rel="noreferrer">
						contributors
					</a>
				</span>
				<a href={githubUrl} target="_blank" rel="noreferrer">
					<StarIcon size={16} weight="bold" aria-hidden="true" />
					Give us a star on GitHub
					<GithubLogoIcon size={16} weight="bold" aria-hidden="true" />
				</a>
			</div>
		</footer>
	);
}
