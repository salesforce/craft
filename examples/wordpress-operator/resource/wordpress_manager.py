import os
import re
import ast
import sys
import yaml
import json
import argparse
import subprocess
from jinja2 import Environment, FileSystemLoader
import glob
import io
import time
import hvac
import base64

base_dir =  os.getcwd()

# def get_secret(secret, path, key):

#     env = os.environ['environment']
#     if os.environ['environment'] == "falcon":

#         vault_addr = os.environ['vault_addr']

#         awsclient = hvac.Client(url=vault_addr, verify=False)

#         namespace = os.environ['namespace']
#         user= os.environ['vault_kubernetes_secret']
#         process = subprocess.Popen(['kubectl', '-n', namespace, 'get', 'secret', user, "-o", "json"], stdout=subprocess.PIPE)
#         process.wait()
#         data, err = process.communicate()
#         if process.returncode is 0:
#             json_data = json.loads(data.decode('utf-8'))

#         decode_data = base64.b64decode(json_data['data']['credentials.json']).decode("utf-8")
#         cred_data = json.loads(decode_data)

#         access_key = cred_data["AccessKeyId"]
#         secret_key = cred_data["SecretAccessKey"]
#         session_token = cred_data["SessionToken"]

#         header_value = "api.vault.secrets." + os.environ['iac_instance'] + ".aws.sfdc.cl"
#         role = "kv_" + os.environ['vault_rootpath'] + "-ro"
#         vault_aws = awsclient.auth.aws.iam_login(access_key=access_key,
#                                 secret_key=secret_key,session_token=session_token,
#                                 header_value=header_value,
#                                 role=role)

#         client = hvac.Client(
#             url=vault_addr,
#             token=vault_aws['auth']['client_token'],
#             verify=False
#         )

       
#         return str(client.read(secret + '/data/' + path)['data'][key])
#     else:
#         client = hvac.Client(
#             url=os.environ['vault_addr'],
#             token=os.environ['VAULT_TOKEN']
#         )

#         return str(client.read(secret + '/' + path)['data'][key])

# def decrypt_passwords(spec):

#     spec_t = {}
#     for data in spec:
#         if str(spec[data]).startswith('vault:'):
#             t_vault = spec[data].split(":")
#             spec_t[data] = get_secret(t_vault[1],t_vault[2],t_vault[3])
#             print(spec_t[data])
#         else:
#             spec_t[data] = spec[data]
            
#     return spec_t

def init_wordpress(spec):

    pods_list = []
    while True:
        retry = False
        time.sleep(10)
        result = subprocess.Popen(['kubectl', 'get', 'pods', '-l=app=wordpress-' + spec['instance']], stdout=subprocess.PIPE)
        for line in io.TextIOWrapper(result.stdout, encoding="utf-8"):
            if 'wordpress-' + spec['instance'] in line:
                if 'Running' not in line:
                    time.sleep(5)
                    retry = True
                else:
                    pods_list.append(line.split(' ')[0])

        if not retry:
            break

    time.sleep(20)
    for pod in pods_list:
        result = subprocess.run(['kubectl', 'exec', '-it', pod, '--', 'sh', '-x', '/scripts/initwordpress.sh', spec['bootstrap_title'], spec['bootstrap_user'], spec['bootstrap_password'],  spec['bootstrap_email'] , spec['instance'], spec['bootstrap_url']], stdout=subprocess.PIPE)
        print(result)
    
def create_init(spec):
    result = subprocess.run(['kubectl', 'create', 'configmap', 'initwordpress-' + spec['instance'], '--from-file=initwordpress.sh'], stdout=subprocess.PIPE)
    print(result)
    return

def delete_init(spec):
    result = subprocess.run(['kubectl', 'delete', 'configmap', 'initwordpress-' + spec['instance']], stdout=subprocess.PIPE)
    print(result)
    return

def create_wordpress(spec):

    file_loader = FileSystemLoader('')
    env = Environment(loader=file_loader)
    env.trim_blocks = True
    env.lstrip_blocks = True
    env.rstrip_blocks = True

    print("Create")
    # spec = decrypt_passwords(spec)
    create_init(spec)
    try:
        os.mkdir( spec['instance'] )
    except:
        print("Dir exists")
    for template in  ['templates/*.yaml']:
        for filename in glob.iglob(template, recursive=True):
            print(filename)
            template = env.get_template( filename )
            new_filename = template
            head, new_filename = os.path.split(filename)
            _ = head

            output = template.render(instance=spec['instance'], replicas=spec['replicas'], db_password=spec['db_password'], dbVolumeMount=spec['dbVolumeMount'], wordpressVolumeMount=spec['wordpressVolumeMount'])
            newpath = os.path.join( spec['instance'] + '/' + new_filename)
            with open(newpath, 'w') as f:
                f.write(output)

    result = subprocess.run(['kubectl', 'apply', '-k', spec['instance']], stdout=subprocess.PIPE)
    print(result)
    print("***********************")

    if result.returncode == 0:
        init_wordpress(spec)
        sys.exit(201)
    else:
        sys.exit(203)


def delete_wordpress(spec):

    print("Delete")
    file_loader = FileSystemLoader('')
    env = Environment(loader=file_loader)
    env.trim_blocks = True
    env.lstrip_blocks = True
    env.rstrip_blocks = True

    # spec = decrypt_passwords(spec)
    delete_init(spec)
    try:
        os.mkdir( spec['instance'] )
    except:
        print("Dir exists")
    for template in  ['templates/*.yaml']:
        for filename in glob.iglob(template, recursive=True):
            print(filename)
            template = env.get_template( filename )
            new_filename = template
            head, new_filename = os.path.split(filename)
            _ = head

            output = template.render(instance=spec['instance'], replicas=spec['replicas'], db_password=spec['db_password'], dbVolumeMount=spec['dbVolumeMount'], wordpressVolumeMount=spec['wordpressVolumeMount'])
            newpath = os.path.join( spec['instance'] + '/' + new_filename)
            with open(newpath, 'w') as f:
                f.write(output)

    result = subprocess.run(['kubectl', 'delete', '-k', spec['instance']], stdout=subprocess.PIPE)
    print(result)
    print("***********************")
    if result.returncode == 0:
        sys.exit(221)
    else:
        sys.exit(223)

def update_wordpress(spec):

    print("Update")
    file_loader = FileSystemLoader('')
    env = Environment(loader=file_loader)
    env.trim_blocks = True
    env.lstrip_blocks = True
    env.rstrip_blocks = True
    # spec = decrypt_passwords(spec)
    try:
        os.mkdir( spec['instance'] )
    except:
        print("Dir exists")
    for template in  ['templates/*.yaml']:
        for filename in glob.iglob(template, recursive=True):
            print(filename)
            template = env.get_template( filename )
            new_filename = template
            head, new_filename = os.path.split(filename)
            _ = head

            output = template.render(instance=spec['instance'], replicas=spec['replicas'], db_password=spec['db_password'], dbVolumeMount=spec['dbVolumeMount'], wordpressVolumeMount=spec['wordpressVolumeMount'])

            newpath = os.path.join( spec['instance'] + '/' + new_filename)
            with open(newpath, 'w') as f:
                f.write(output)

    result = subprocess.run(['kubectl', 'delete', '-k', spec['instance']], stdout=subprocess.PIPE)
    print(result)
    if result.returncode == 0:
        result = subprocess.run(['kubectl', 'apply', '-k', spec['instance']], stdout=subprocess.PIPE)
        print(result)
        if result.returncode == 0:
            sys.exit(201)
        else:
            sys.exit(203)

def verify_wordpress(spec):

    print("Verify")
    result = subprocess.run(['kubectl', 'get', 'deployment', 'wordpress-' + spec['instance'] ], stdout=subprocess.PIPE)
    print(result)
    print("***********************")
    if result.returncode == 0:
        result = subprocess.run(['kubectl', 'get', 'deployment', 'wordpress-' + spec['instance'],  '-o', 'yaml'], stdout=subprocess.PIPE)
        deployment_out = yaml.safe_load(result.stdout)
        if deployment_out['spec']['replicas'] != spec['replicas']: 
            print("Change in replicas.")
            sys.exit(214)
        sys.exit(211)
    else:
        sys.exit(214)


def convert_to_dict(tmp_config):
    t_config=tmp_config
    if not isinstance(t_config,dict):
        try:
            t_config = t_config.replace('\'','\"')
            t_config=json.loads(t_config)
        except:
            t_config=ast.literal_eval(tmp_config)
        return t_config
    return tmp_config

if __name__ == "__main__":

    parser = argparse.ArgumentParser(description='Initiating clone_autobuild')
    parser.add_argument('--type', type=str, required=True)
    parser.add_argument('--spec', type=str, required=True)

    args = parser.parse_args()
    action_type = args.type
    # print(args.spec)
    input_data = convert_to_dict(args.spec)

    if action_type == 'create':
        create_wordpress(input_data)
    if action_type == 'verify':
        verify_wordpress(input_data)
    if action_type == 'update':
        update_wordpress(input_data)
    if action_type == 'delete':
        delete_wordpress(input_data)

