import * as React from "react";
import { Link } from "react-router-dom";

export interface NavigationProps {
	brandName: string;
}

export class Navigation extends React.Component<NavigationProps, {}> {
	render() {
		return <nav className="navbar navbar-default">
			<div className="container">
				<div className="navbar-header">
					<Link to="/" className="navbar-brand">{this.props.brandName}</Link>
				</div>

				<ul className="nav navbar-nav">
					<li className="active"><Link to="/">Events</Link></li>
				</ul>
			</div>
		</nav>
	}
}
