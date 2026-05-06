import { startTransition, useEffect, useState } from "react";
import type { Route } from "../utils/mockData";

const fallbackRoute: Route = "landing";

function isRoute(value: string): value is Route {
	return ["landing", "login", "register", "onboarding", "dashboard", "probes", "insight", "checks", "alerts", "team", "settings", "components"].includes(value);
}

function readHash(): Route {
	const value = window.location.hash.replace("#", "");
	return isRoute(value) ? value : fallbackRoute;
}

export function useHashNavigation() {
	const [route, setRoute] = useState(() => readHash());

	useEffect(() => {
		function handleHashChange() {
			startTransition(() => setRoute(readHash()));
		}

		window.addEventListener("hashchange", handleHashChange);
		return () => window.removeEventListener("hashchange", handleHashChange);
	}, []);

	function navigate(nextRoute: Route) {
		if (nextRoute === route) {
			return;
		}

		window.location.hash = nextRoute;
		startTransition(() => setRoute(nextRoute));
	}

	return { route, navigate };
}
