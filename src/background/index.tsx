import {Parallax} from 'react-scroll-parallax';

export function ScrollingBackground(){
	return(
		<div class="static">
		<Parallax speed={-10}>
		<img src="/og-photo.jpg"></img>
		<img src="/og-photo.jpg"></img>
		<img src="/og-photo.jpg"></img>
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
