# FeatureSpace Fraud Analytics - Scenario post'er

Little app to post various sets of JSON files (from $input_path) to build up a test scenario to test rules/alerts

Create a certs directory and copy client.crt and client.key into it.

1. Create a json_output directory

2. Create a json_source# directory for every set of scenario sets to be tested.

3. Enable/Disable options in .exps

4. Modify the .exps to point to the relevant json_source# directory

5. execute .exps to configure environment variables.

6. execute fs_producer.exe