import React, { Component } from "react";
import client from "./client";
import AftTable from "./components/AftTable";
import AftNav from "./components/AftNav";
import AftForm from "./components/AftForm";
import AftError from "./components/AftError";

/* the main page for the index route of this app */
class App extends Component {
  constructor() {
    super();
    this.state = { models: [], new: false, selected: "", error: "" };
    this.handleSelect = this.handleSelect.bind(this);
    this.handleNew = this.handleNew.bind(this);
    this.handleError = this.handleError.bind(this);
    let load = client.api.model.findMany({ where: { system: false } });
    load.then(obj => {
      var models = [];
      obj.map(val => models.push(val["name"]));
      if (models.length > 0) {
        this.setState({ models: models, new: false, selected: models[0] });
      }
    });
  }

  handleSelect(model) {
    this.setState(prevState => {
      return {
        ...prevState,
        selected: model,
        new: false
      };
    });
  }

  handleNew() {
    this.setState(prevState => {
      return {
        ...prevState,
        new: true
      };
    });
  }

  handleError(e) {
    console.log(e);
    this.setState(prevState => {
      return {
        ...prevState,
        error: e
      };
    });
  }

  render() {
    if (this.state.selected !== "") {
      var main = "";
      if (this.state.new) {
        main = (
          <AftForm
            name={this.state.selected}
            handleSubmit={this.handleSelect}
            handleError={this.handleError}
          />
        );
      } else {
        main = <AftTable name={this.state.selected} />;
      }

      return (
        <div style={{ padding: "1em 1.5em" }}>
          <AftError error={this.state.error} />
          <AftNav
            name={this.state.selected}
            list={this.state.models}
            handleSelect={this.handleSelect}
            handleNew={this.handleNew}
          />
          {main}
        </div>
      );
    }
    return <div>No models to play with. Go to aft to create some...</div>;
  }
}
export default App;
