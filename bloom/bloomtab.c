/* Copyright (c) 2011-2014 PLUMgrid, http://plumgrid.com
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of version 2 of the GNU General Public
 * License as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 */
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/bpf.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Tristan Colgate-McFarlane<tcolgate@gmail.com>");
MODULE_DESCRIPTION("Bloom filter table support for eBPF");

struct bpf_bloom {
	struct bpf_map map;
	spinlock_t lock;
	u32 count;	/* number of elements in this hasbloomle */
	u32 n_bits;	/* number of hash buckets */
};

/* Called from syscall */
static struct bpf_map *bloom_map_alloc(union bpf_attr *attr)
{
        printk("* In function %s *\n", __FUNCTION__);
	return ERR_PTR(-ENOMEM);
#if 0
	struct bpf_bloom *bloom;
	int err, i;

	bloom = kzalloc(sizeof(*bloom), GFP_USER);
	if (!bloom)
		return ERR_PTR(-ENOMEM);

	/* mandatory map attributes */
	bloom->map.key_size = attr->key_size;
	bloom->map.value_size = attr->value_size;
	bloom->map.max_entries = attr->max_entries;

	/* check sanity of attributes.
	 * value_size == 0 may be allowed in the future to use map as a set
	 */
	err = -EINVAL;
	if (bloom->map.max_entries == 0 || bloom->map.key_size == 0 ||
	    bloom->map.value_size == 0)
		goto free_bloom;

	/* hash table size must be power of 2 */
	bloom->n_buckets = roundup_pow_of_two(bloom->map.max_entries);

	err = -E2BIG;
	if (bloom->map.key_size > MAX_BPF_STACK)
		/* eBPF programs initialize keys on stack, so they cannot be
		 * larger than max stack size
		 */
		goto free_bloom;

	err = -ENOMEM;
	/* prevent zero size kmalloc and check for u32 overflow */
	if (bloom->n_buckets == 0 ||
	    bloom->n_buckets > U32_MAX / sizeof(struct hlist_head))
		goto free_bloom;

	bloom->buckets = kmalloc_array(bloom->n_buckets, sizeof(struct hlist_head),
				      GFP_USER | __GFP_NOWARN);

	if (!bloom->buckets) {
		bloom->buckets = vmalloc(bloom->n_buckets * sizeof(struct hlist_head));
		if (!bloom->buckets)
			goto free_bloom;
	}

	for (i = 0; i < bloom->n_buckets; i++)
		INIT_HLIST_HEAD(&bloom->buckets[i]);

	spin_lock_init(&bloom->lock);
	bloom->count = 0;

	bloom->elem_size = sizeof(struct bloom_elem) +
			  round_up(bloom->map.key_size, 8) +
			  bloom->map.value_size;
	return &bloom->map;

free_bloom:
	kfree(bloom);
	return ERR_PTR(err);
#endif
}

/* Called from syscall or from eBPF program */
static void *bloom_map_lookup_elem(struct bpf_map *map, void *key)
{
        printk("* In function %s *\n", __FUNCTION__);
	return NULL;
}

/* Called from syscall */
static int bloom_map_get_next_key(struct bpf_map *map, void *key, void *next_key)
{
        printk("* In function %s *\n", __FUNCTION__);
	return -ENOENT;
}

/* Called from syscall or from eBPF program */
static int bloom_map_update_elem(struct bpf_map *map, void *key, void *value,
				u64 map_flags)
{
        printk("* In function %s *\n", __FUNCTION__);
  return -ENOENT;
}

/* Called from syscall or from eBPF program */
static int bloom_map_delete_elem(struct bpf_map *map, void *key)
{
        printk("* In function %s *\n", __FUNCTION__);
  return -ENOENT;
}

static void delete_all_elements(struct bpf_bloom *bloom)
{
}

/* Called when map->refcnt goes to zero, either from workqueue or from syscall */
static void bloom_map_free(struct bpf_map *map)
{
        printk("* In function %s *\n", __FUNCTION__);
}

static const struct bpf_map_ops bloom_ops = {
	.map_alloc = bloom_map_alloc,
	.map_free = bloom_map_free,
	.map_get_next_key = bloom_map_get_next_key,
	.map_lookup_elem = bloom_map_lookup_elem,
	.map_update_elem = bloom_map_update_elem,
	.map_delete_elem = bloom_map_delete_elem,
};

static struct bpf_map_type_list bloom_type __read_mostly = {
	.ops = &bloom_ops,
//	.type = BPF_MAP_TYPE_BLOOM,
	.type = BPF_MAP_TYPE_UNSPEC,
};

static int __init register_bloom_map(void)
{
        printk("* In function %s *\n", __FUNCTION__);
	bpf_register_map_type(&bloom_type);
	return 0;
}

static void __exit unregister_bloom_map(void)
{
        printk("* In function %s *\n", __FUNCTION__);
	return;
}

//late_initcall(register_bloom_map);

module_init(register_bloom_map);
module_exit(unregister_bloom_map);



