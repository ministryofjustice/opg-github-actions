# Timer Composite Action

Provides simple timer functionality so you can track how long sections of workflows take to run - for exmaple image builds or infrastructure deployments

## Usage

Within you github workflow job you can place this to start a timer:

```yaml
    - name: "Start build timer"
      id: build_start
      uses: 'ministryofjustice/opg-github-actions/.github/actions/timer@v3.0.7'
```

Then, later in your workflow you can stop that timer and get duration information by including the below style step:

```yaml
    - name: "Stop build timer"
      id: build_stop
      uses: 'ministryofjustice/opg-github-actions/.github/actions/timer@v3.0.7'
      with:
        stop: true
        timestamp: ${{ steps.build_start.outputs.start }}
```

The `timestamp` variable passed is a IS0 8601 timestamp (UTC timezone) in the below format:

```
2024-03-26T08:47:21.029473+00:00
```

*Please note the `+00:00` usage rather than `Z` for timezone offsets.*


## Inputs and Outputs

Inputs:
- `start` (default: "true")
- `stop`
- `timestamp` 

Outputs:
- **`duration`**
- `start`
- `end`
- `duration_as_milliseconds`
- `duration_as_seconds`
- `duration_as_minutes`
- `duration_as_hours`


### Inputs

#### `start` (default: "true")
While this is true and `stop` has no value, then this will trigger the generation of a timestamp

#### `stop` and `timestamp` 
When `stop` is set to true and a timestamp value is passed this will trigger the output of duration information between the current time and the value of timestamp.


### Outputs

#### `duration`
The difference in seconds between the passed timestamp and the current time. This is rounded to the nearest second, for more accuracy you can use `duration_as_milliseconds` or `duration_as_seconds`.

#### `start` and `end`
The two timestamp values used to calculate the duration data.

#### `duration_as_milliseconds`, `duration_as_seconds`, `duration_as_minutes`, `duration_as_hours`
Additional duration formats that provide deciminal versions of there respective unit of time.


