import React, { Component } from "react";
import client from "../client";
import Form from "@rjsf/core";
import { isObject } from "../util";

class AftForm extends Component {
  constructor(props) {
    super(props);
    this.state = {};
    this.handleSubmit = this.handleSubmit.bind(this);

    if (this.props.name !== "") {
      let load = client.rpc.reactForm({
        data: { model: this.props.name }
      });
      load.then(obj => {
        this.setState({ schema: obj["schema"], uiSchema: obj["uiSchema"] });
      });
    }
  }

  handleSubmit(e) {
    this.setState(prevState => {
      return {
        ...prevState,
        formData: e.formData
      };
    });
    let submit = client.api[this.props.name].create({ data: e.formData });
    submit.then(obj => {
      this.props.handleSubmit(this.props.name);
    });

    let errors = client.rpc.validate({
      data: { schema: this.state.schema, data: e.formData }
    });
    errors.then(obj => {
      this.setState(prevState => {
        return {
          ...prevState,
          errors: obj
        };
      });
    });
  }

  render() {
    if (isObject(this.state.schema)) {
      return (
        <Form
          schema={this.state.schema}
          onSubmit={fd => this.handleSubmit(fd)}
          extraErrors={this.state.errors}
          formData={this.state.formData}
          uiSchema={this.state.uiSchema}
        />
      );
    }
    return <div></div>;
  }
}

export default AftForm;
