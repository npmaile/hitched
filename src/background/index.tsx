import React from 'react';
import { Parallax } from 'react-scroll-parallax';

export function ScrollingBackground() {
	return (
		<div className="absolute z-0">
			<Parallax speed={100}>
				<img src="/og-photo.jpg"></img>
				<img src="/og-photo.jpg"></img>
				<img src="/og-photo.jpg"></img>
				<img src="/og-photo.jpg"></img>
				<img src="/og-photo.jpg"></img>
				<img src="/og-photo.jpg"></img>
			</Parallax>
		</div>
	);
}
