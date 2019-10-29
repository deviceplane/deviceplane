import React, { Fragment, Component } from "react";
import {
  Pane,
  minorScale,
} from "evergreen-ui";

export function buildLabelColorMap(oldLabelColorMap, labelColors, items) {
  var x = [];
  items.forEach(item => {
    Object.keys(item.labels).forEach(label => {
      x.push(label);
    });
  });
  const labelKeys = [...new Set(x)].sort();

  var labelColorMap = Object.assign({}, oldLabelColorMap);
  labelKeys.forEach((key, i) => {
    if (!labelColorMap[key]) {
      labelColorMap[key] = labelColors[i % (labelColors.length - 1)]
    }
  });
  return labelColorMap;
};

export function renderLabels(labels, labelColorMap) {
  return (
    <Pane display="flex" flexWrap="wrap">
      {Object.keys(labels).map((key, i) => (
        <Pane
          display="flex"
          marginRight={minorScale(2)}
          marginY={minorScale(1)}
          overflow="hidden"
          key={key}
        >
          <Pane
            backgroundColor={labelColorMap[key]}
            paddingX={6}
            paddingY={2}
            color="white"
            borderTopLeftRadius={3}
            borderBottomLeftRadius={3}
            textOverflow="ellipsis"
            overflow="hidden"
            whiteSpace="nowrap"
          >
            {key}
          </Pane>
          <Pane
            backgroundColor="#E4E7EB"
            paddingX={6}
            paddingY={2}
            borderTopRightRadius={3}
            borderBottomRightRadius={3}
            textOverflow="ellipsis"
            overflow="hidden"
            whiteSpace="nowrap"
          >
            {labels[key]}
          </Pane>
        </Pane>
      ))}
    </Pane>
  );
}