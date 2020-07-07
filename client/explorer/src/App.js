import React, { Component } from "react";
import client from "./client";
import AftTable from "./components/AftTable";
import AftNav from "./components/AftNav";
import AftForm from "./components/AftForm";
import "./components/aft.css";
import { cap } from "./util";

/* the main page for the index route of this app */
class App extends Component {
  constructor() {
    super();
    this.state = { models: [], new: false, selected: "" };
    this.handleNav = this.handleNav.bind(this);
    this.handleSubmitModel = this.handleSubmitModel.bind(this);
    this.handleNewModel = this.handleNewModel.bind(this);
    let load = client.api.model.findMany({ where: { system: false } });
    load.then(obj => {
      var models = [];
      obj.map(val => models.push(val["name"]));
      if (models.length > 0) {
        this.setState({ models: models, new: false, selected: models[0] });
      }
    });
  }

  handleSubmitModel(e) {
    this.setState(prevState => {
      return {
        ...prevState,
        new: false
      };
    });
  }

  handleNav(model) {
    this.setState(prevState => {
      return {
        ...prevState,
        selected: model,
        new: false
      };
    });
  }
  
  handleNewModel() {
    this.setState(prevState => {
      return {
        ...prevState,
        new: true
      };
    });
  }

  render() {
    if (this.state.selected !== "") {
      var main = "";
      if (this.state.new) {
        main = (
          <AftForm
            model={this.state.selected}
            handleSuccess={this.handleSubmitModel}

            title={cap(this.state.selected)}
            titleClassName="text-primary h2 font-weight-light"
            fieldClassName="field-pad"
            buttonText={"Submit"}
            buttonClassName="btn btn-light"
            labelClassName="text-primary font-weight-light font-italic text-capitalize"
            className="form text-primary bg-secondary"
          />
        );
      } else {
        main = <AftTable name={this.state.selected} />;
      }

      return (
        <div style={{ padding: "1em 1.5em" }}>
          <AftNav
            name={this.state.selected}
            list={this.state.models}
            handleSelect={this.handleNav}
            handleNew={this.handleNewModel}
          />
          {main}
        </div>
      );
    }
    return <div>No models to play with. Go to aft to create some...</div>;
  }
}

export default App;
