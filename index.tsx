import { hydrateRoot } from "react-dom/client";
import { Component } from "react";
document.body.innerHTML = '<div id="app"></div>';
import imgUrl from "./photos/og-photo.jpg?w=300&h=300";
document.getElementById("meta-image")!.setAttribute("content", imgUrl);

class Abc extends Component {
  render() {
    return <div>information will come soon</div>;
  }
}

const root = hydrateRoot(
  document.getElementById("app"),
  <h1>
    Nate and Susy wedding page
    <Abc />
    <img src={imgUrl}></img>
    <a href="https://www.amazon.com/wedding/share/SN8HITCHED">registry</a>
  </h1>,
);
