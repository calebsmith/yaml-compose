Noop - Output for a normal YAML file should be invariant:
  $ pushd $INPUTS > /dev/null; yaml-compose vars.yaml; popd > /dev/null
  A: 1
  B: 2
  data:
    - C: 3
    - D: 4

Inject - Only performing injection (inlining a single pair):
  $ pushd $INPUTS > /dev/null; yaml-compose inject.yaml; popd > /dev/null
  data:
    keya1: val1
    keya2: val2
    var3key: value

Load - Load a single variable and render it:
  $ pushd $INPUTS > /dev/null; yaml-compose load.yaml; popd > /dev/null
  Value: 1

Full:
  $ pushd $INPUTS > /dev/null; yaml-compose full.yaml; popd > /dev/null
  data:
    - nested:
        keya1: val1
        keya2: val2
  loaded_scalars:
    - 1
    - 2
    - 3
    - 4
  loaded_hashmap:
    key1: value1
    key2: value2
